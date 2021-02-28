package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"sort"
	"strings"
)

const rootDir = filename("assets")

const enforceLeadingDigits = true

var fileNameFormat = regexp.MustCompile("(^\\d+[-].+)")

var endpoints = Endpoints{
	"Dashboards":              "/api/config/v1/dashboards",
	"DashboardSharingDetails": "/api/config/v1/dashboards/{id}/shareSettings",
	"MetricsService":          "/api/config/v1/calculatedMetrics/service",
	"Application":             "/api/config/v1/applications/web",
	"DetectionRules":          "/api/config/v1/applicationDetectionRules",
	"ManagementZones":         "/api/config/v1/managementZones",
	"CalculatedMetrics":       "/api/config/v1/calculatedMetrics/service",
	"CustomServices":          "/api/config/v1/service/customServices/java?position=APPEND",
	"Extensions":              "/api/config/v1/extensions",
	"AutoTags":                "/api/config/v1/autoTags",
	"RequestAttributes":       "/api/config/v1/service/requestAttributes",
}

func main() {
	var err error
	var cookieJar *cookiejar.Jar
	if cookieJar, err = cookiejar.New(nil); err != nil {
		panic(err)
	}
	client := &http.Client{Jar: cookieJar}
	processor := &Processor{
		Config: new(Config).Parse(),
		Vars:   make(variables),
		Client: client,
	}
	if err = processor.Process(); err != nil {
		panic(err)
	}
}

// Processor has no documentation
type Processor struct {
	Config *Config
	Vars   variables
	Client *http.Client
}

// Process has no documentation
func (p *Processor) Process() error {
	var err error
	var folders []filename
	if folders, err = p.readDir(string(rootDir), endpoints.Contains); err != nil {
		return err
	}

	// Perform an early sanity check without sending anything
	// Benefit: If some of the folders or files are illegal/flawed nothing has been sent yet
	for _, dir := range folders {
		if _, err = p.readDir(rootDir.Join(dir), isSupportedFile); err != nil {
			return err
		}
		if _, found := endpoints[dir.WithoutPrefix()]; !found {
			return fmt.Errorf("Endpoint URL not found for %v", dir)
		}
	}

	for _, dir := range folders {
		var files []filename
		if files, err = p.readDir(rootDir.Join(dir), isSupportedFile); err != nil {
			return err
		}
		// Extract the directory name without the number
		endpointName := dir.WithoutPrefix()
		var endpoint string
		var found bool
		if endpoint, found = endpoints[endpointName]; !found {
			return fmt.Errorf("Endpoint URL not found for %v", dir)
		}

		for _, file := range files {
			endpointURL := p.Config.URL + endpoint
			if p.Config.Verbose {
				log.Println("==== " + rootDir.Join("/", dir, file) + " ====")
			}
			if endpointName == "Extensions" {
				if err = p.sendExtensionFile(endpointURL, rootDir.Join(dir, file)); err != nil {
					return err
				}
				continue
			}

			var data []byte
			if data, err = p.Vars.Load(rootDir.Join(dir, file)); err != nil {
				return err
			}
			var id string
			if id, err = p.post(endpointName, endpointURL, data); err != nil {
				return err
			}

			if endpointName == "Dashboards" {
				if p.Config.Verbose {
					log.Println("Sending sharing data")
				}
				if err := p.sendSharingDetails(id); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func (p *Processor) post(endpointName string, endpointURL string, data []byte) (string, error) {
	var err error
	var req *http.Request
	if req, err = p.setupHTTPRequest("POST", endpointURL, data); err != nil {
		return "", err
	}
	var resp *http.Response
	if resp, err = p.Client.Do(req); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var id string
	id, _, err = p.processResponse(resp, endpointName)
	return id, err
}

func (p *Processor) setupHTTPRequest(method string, endpointURL string, data []byte) (*http.Request, error) {
	if p.Config.Debug {
		log.Println(fmt.Sprintf("  [HTTP] %s %s", method, endpointURL))
		log.Println(fmt.Sprintf("  [HTTP] %s", string(data)))
	}
	var err error
	var req *http.Request
	if req, err = http.NewRequest(method, endpointURL, bytes.NewBuffer(data)); err != nil {
		return nil, err
	}
	if p.Config.Debug {
		log.Println(fmt.Sprintf("  [HTTP] %s: %s", "Content-Type", "application/json"))
	}
	req.Header.Set("Content-Type", "application/json")

	if p.Config.Debug {
		log.Println(fmt.Sprintf("  [HTTP] %s: %s", "Authorization", "Api-Token "+p.Config.APIToken))
	}
	req.Header.Add("Authorization", "Api-Token "+p.Config.APIToken)

	return req, nil
}

func (p *Processor) sendSharingDetails(dashboardID string) error {
	var err error
	endpointURL := endpoints["DashboardSharingDetails"]
	replacedURL := p.Config.URL + strings.Replace(endpointURL, "{id}", dashboardID, 1)
	sharingDetails := SharingDetails{
		ID:           dashboardID,
		Published:    true,
		Enabled:      true,
		Permissions:  []Permission{{Type: "ALL", Permission: "VIEW"}},
		PublicAccess: PublicAccess{ManagementZoneIDs: []string{}},
	}
	var data []byte
	if data, err = json.Marshal(&sharingDetails); err != nil {
		return err
	}
	// data := "{\"id\": \"{id}\",\"enabled\": true,\"published\": true,\"permissions\": [{\"type\": \"ALL\",\"permission\": \"VIEW\"}],\"publicAccess\": {\"managementZoneIds\": [],\"urls\": {}}}"
	// replacedData := strings.Replace(data, "{id}", dashboardID, 1)
	//fmt.Printf("replacedData=%s, replacedURL=%s\n", replacedData, replacedURL)
	var req *http.Request
	if req, err = p.setupHTTPRequest("PUT", replacedURL, []byte(data)); err != nil {
		return err
	}
	var resp *http.Response
	if resp, err = p.Client.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, _, err = p.processResponse(resp, "DashboardSharingDetails"); err != nil {
		return err
	}
	return nil
}

func (p *Processor) processResponse(resp *http.Response, endpointName string) (string, string, error) {
	fmt.Println("processResponse")
	body, _ := ioutil.ReadAll(resp.Body)

	if p.Config.Debug {
		log.Println("  [HTTP] Status: " + resp.Status)
	}
	if resp.StatusCode >= 400 {
		var errorEnvelope ErrorEnvelope
		json.Unmarshal([]byte(body), &errorEnvelope)
		if errorEnvelope.RESTError == nil {
			var restError RESTError
			json.Unmarshal([]byte(body), &restError)
			return "", "", &restError
		}
		return "", "", &errorEnvelope
	}
	var successEnvelope SuccessEnvelope
	json.Unmarshal([]byte(body), &successEnvelope)
	if successEnvelope.ID != "" && p.Config.Verbose {
		log.Println(fmt.Sprintf("  id: %s", successEnvelope.ID))
		if len(successEnvelope.Name) > 0 {
			log.Println(fmt.Sprintf("  name: %s", successEnvelope.Name))
		}
	}
	varKey := fmt.Sprintf("%s.id", strings.ToLower(endpointName))
	if p.Config.Verbose {
		log.Println("  [ENV] {" + varKey + "} <= " + successEnvelope.ID)
	}
	p.Vars[varKey] = successEnvelope.ID
	return successEnvelope.ID, successEnvelope.Name, nil

}

func (p *Processor) upload(url string, values map[string]io.Reader) (*http.Response, error) {
	var err error
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return nil, err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return nil, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	var req *http.Request
	// Now that you have a form, you can submit it to your handler.
	if req, err = http.NewRequest("POST", url, &b); err != nil {
		return nil, err
	}
	if p.Config.Debug {
		log.Println(fmt.Sprintf("  [HTTP] %s: %s", "Content-Type", w.FormDataContentType()))
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	if p.Config.Debug {
		log.Println(fmt.Sprintf("  [HTTP] %s: %s", "Authorization", "Api-Token "+p.Config.APIToken))
	}
	req.Header.Add("Authorization", "Api-Token "+p.Config.APIToken)

	// Submit the request
	return p.Client.Do(req)
}

func (p *Processor) sendExtensionFile(url string, filename string) error {
	var err error
	var file *os.File

	if file, err = os.Open(filename); err != nil {
		return err
	}
	defer file.Close()

	var resp *http.Response
	if resp, err = p.upload(url, map[string]io.Reader{
		"file": file,
	}); err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, _, err = p.processResponse(resp, "Extensions"); err != nil {
		return err
	}
	return nil
}

func (p *Processor) readDir(dirName string, match func(os.FileInfo) bool) ([]filename, error) {
	var dirs []os.FileInfo
	var err error
	if dirs, err = ioutil.ReadDir(dirName); err != nil {
		return nil, err
	}
	fileNames := []string{}
	for _, fileInfo := range dirs {
		if match(fileInfo) {
			fileNames = append(fileNames, fileInfo.Name())
		} else if p.Config.Verbose {
			log.Println(fmt.Sprintf("... ignoring %s/%s", dirName, fileInfo.Name()))
		}
	}
	sort.Strings(fileNames)
	result := make([]filename, len(fileNames))
	for i, v := range fileNames {
		result[i] = filename(v)
	}
	return result, nil
}

/********************* CONFIGURATION *********************/

// Config a simple configuration object
type Config struct {
	URL      string
	APIToken string
	Verbose  bool
	Debug    bool
}

// Parse reads configuration from arguments and environment
func (c *Config) Parse() *Config {
	flag.StringVar(&c.URL, "url", "", "the Dynatrace environment URL (e.g. https://####.live.dynatrace.com)")
	flag.StringVar(&c.APIToken, "token", "", "the API token to use for uploading configuration")
	flag.BoolVar(&c.Verbose, "verbose", false, "verbose logging")
	flag.BoolVar(&c.Debug, "debug", false, "prints out HTTP traffic")
	flag.Parse()
	c.URL = c.Lookup("DT_URL", c.URL)
	c.APIToken = c.Lookup("DT_TOKEN", c.APIToken)
	if len(c.URL) == 0 || len(c.APIToken) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	return c
}

// Lookup reads configuration from environment
func (c *Config) Lookup(envVar string, current string) string {
	if len(current) > 0 {
		return current
	}
	if v, found := os.LookupEnv(envVar); found {
		return v
	}
	return current
}

func isSupportedFile(fileInfo os.FileInfo) bool {
	if fileInfo.IsDir() {
		return false
	}
	name := fileInfo.Name()
	return strings.HasSuffix(name, ".json") || strings.HasSuffix(name, ".zip")
}

/********************* VARIABLE SUBSTITUTION *********************/
type variables map[string]string

// Load reads the given JSON file from disk, and automatically substitutes
// potential variables with the currently stored values.
func (vars variables) Load(filePath string) ([]byte, error) {
	var err error
	var data []byte
	if data, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}

	model := map[string]interface{}{}

	if err = json.Unmarshal(data, &model); err != nil {
		return nil, err
	}
	vars.replace(model)
	if data, err = json.Marshal(model); err != nil {
		return nil, err
	}
	return data, nil
}

func (vars variables) replace(o interface{}) interface{} {
	if o == nil {
		return nil
	}
	switch to := o.(type) {
	case string:
		if strings.HasPrefix(to, "{") && strings.HasSuffix(to, "}") {
			key := to[1 : len(to)-1]
			if value, found := vars[key]; found {
				return value
			}
		}
		return to
	case []interface{}:
		for i, v := range to {
			to[i] = vars.replace(v)
		}
		return to
	case map[string]interface{}:
		for k, v := range to {
			to[k] = vars.replace(v)
		}
		return to
	case int, int8, int32, int16, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return to
	default:
		panic(fmt.Sprintf("unsupported: %t", o))
	}
}

/********************* CONVENIENCE TYPES *********************/

// filename is nothing else than a string
// It can produce the file name without prefix
// and it can ajoin additional files (file path)
type filename string

func (fn filename) String() string {
	return string(fn)
}

func (fn filename) WithoutPrefix() string {
	s := fn.String()
	return s[strings.Index(s, "-")+1:]
}

func (fn filename) Join(items ...interface{}) string {
	result := fn.String()
	for _, item := range items {
		result = fmt.Sprintf("%v%s%v", result, "/", item)
	}
	return result
}

// Endpoints is a convenience type for a map[string]string
// You can ask this object whether a specific file matches
// the prerequisites for the currently supported endpoint categories
type Endpoints map[string]string

// Contains returns true if an entry exist for this key, false otherwise
func (eps Endpoints) Contains(fileInfo os.FileInfo) bool {
	if fileInfo == nil || !fileInfo.IsDir() {
		return false
	}

	name := fileInfo.Name()

	if enforceLeadingDigits && !fileNameFormat.MatchString((name)) {
		return false
	}
	idx := strings.Index(name, "-")
	if idx >= 0 {
		name = name[idx+1:]
	}
	_, found := eps[name]
	return found
}

// // RestClient has no documentation
// type RestClient struct {
// 	HTTPClient *http.Client
// 	APIToken   string
// 	Debug      bool
// }

// func (rc *RestClient) addHeader(req *http.Request, name string, value string) {
// 	if rc.Debug {
// 		log.Println(fmt.Sprintf("  [HTTP] %s: %s", name, value))
// 	}
// 	req.Header.Add(name, value)
// }

// // Send has no documentation
// func (rc *RestClient) Send(method string, url string, headers map[string]string, data io.Reader) (*http.Response, error) {
// 	var err error
// 	var req *http.Request

// 	if req, err = http.NewRequest(method, url, data); err != nil {
// 		return nil, err
// 	}
// 	rc.addHeader(req, "Authorization", "Api-Token "+rc.APIToken)
// 	if headers != nil {
// 		for name, value := range headers {
// 			rc.addHeader(req, name, value)
// 		}
// 	}

// 	var resp *http.Response

// 	if resp, err = rc.HTTPClient.Do(req); err != nil {
// 		return nil, err
// 	}

// 	return resp, err
// }

/********************* API PAYLOAD *********************/

// ErrorEnvelope is potentially the JSON response code
// when a REST API call fails
type ErrorEnvelope struct {
	RESTError *RESTError `json:"error"` // the actual error object
}

func (e *ErrorEnvelope) Error() string {
	bytes, _ := json.MarshalIndent(e.RESTError, "", "  ")
	return string(bytes)
}

// RESTError is potentially the JSON response code
// when a REST API call fails
type RESTError struct {
	Code                 int                    `json:"code"`    // error code
	Message              string                 `json:"message"` // error message
	ConstraintViolations []*ConstraintViolation `json:"constraintViolations"`
}

func (e *RESTError) Error() string {
	bytes, _ := json.MarshalIndent(e, "", "  ")
	return string(bytes)
}

// ConstraintViolation is the viloation error
type ConstraintViolation struct {
	Path              string `json:"path"`
	Message           string `json:"message"`
	ParameterLocation string `json:"parameterLocation"`
	Location          string `json:"location"`
}

// SuccessEnvelope contains the successful response from the API endpoint
type SuccessEnvelope struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PublicAccess has no documentation
type PublicAccess struct {
	ManagementZoneIDs []string `json:"managementZoneIds"`
}

// Permission has no documentation
type Permission struct {
	Type       string `json:"type"`
	Permission string `json:"permission"`
}

// SharingDetails has no documentation
type SharingDetails struct {
	ID           string       `json:"id"`
	Description  string       `json:"description"`
	PublicAccess PublicAccess `json:"publicAccess"`
	Enabled      bool         `json:"enabled"`
	Published    bool         `json:"published"`
	Permissions  []Permission `json:"permissions"`
}
