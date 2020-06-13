package http

//
//
//var (
//	errHttpClientInvalidMethod = errors.New("http request has invalid method")
//)
//var methodsMap = map[string]string{
//	"post":    "POST",
//	"get":     "GET",
//	"head":    "HEAD",
//	"put":     "PUT",
//	"delete":  "DELETE",
//	"patch":   "PATCH",
//	"options": "OPTIONS",
//}
//
//type Client struct {
//	name   string
//	client *resty.Client
//	opts   options
//	log    *logger.Logger
//}
//
//func New() *Client {
//	return &Client{}
//}
//
//func (c *Client) Name() string {
//	return c.name
//}
//func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
//	c.name = cfg.Name
//	c.log = logger.NewLogger(cfg.Name)
//	var err error
//	c.opts, err = parseOptions(cfg)
//	if err != nil {
//		return err
//	}
//	c.client = resty.New()
//	return nil
//}
//func (c *Client) newBaseRequest() *resty.Request {
//	req := c.client.NewRequest()
//	switch c.opts.authType {
//	case "basic":
//		req.SetBasicAuth(c.opts.username, c.opts.password)
//	case "auth_token":
//		req.SetAuthToken(c.opts.token)
//
//	}
//	req.SetHeaders(c.opts.headers)
//	return req
//}
//func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
//	r := c.newBaseRequest()
//	r.URL = fmt.Sprintf("%s%s", c.opts.uri, req.Url)
//	var ok bool
//	r.Method, ok = methodsMap[strings.ToLower(req.Method)]
//	if !ok {
//		c.log.Errorf(errHttpClientInvalidMethod.Error())
//		return nil, errHttpClientInvalidMethod
//	}
//	r.SetContext(ctx)
//	for header, value := range req.Headers {
//		r.SetHeader(header, value)
//	}
//	r.SetDoNotParseResponse(true)
//	r.SetBody(req.Data)
//	resp, err := r.Send()
//	if err != nil {
//		return nil, err
//	}
//	tr, err := newResultFromHttpResponse(resp.RawResponse)
//	if err != nil {
//		c.log.Errorf(errHttpClientInvalidMethod.Error())
//		return nil, err
//	}
//	return tr, resp.RawResponse.Body.Close()
//}
//
//func newResultFromHttpResponse(hr *http.Response) (*types.Response, error) {
//	resp := types.NewResponse().SetCode(hr.StatusCode)
//
//	if hr.StatusCode >= 400 {
//		resp.SetError(hr.Status)
//	}
//
//	for name, values := range hr.Header {
//		resp.Headers[name] = strings.Join(values, ",")
//	}
//	var err error
//	if hr.Body == nil {
//		return resp, nil
//	}
//	resp.Data, err = ioutil.ReadAll(hr.Body)
//	if err != nil {
//		return nil, err
//	}
//	return resp, nil
//}
