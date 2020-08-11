package rpt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang/gddo/httputil/header"
)

// APIServer struct
type APIServer struct {
	BasePath           string
	ListenAddr         string
	Operations         chan *DBOperationSet
	lookupOperationSet map[string]*DBOperationSet
	primary            DBClient
	secondary          DBClient
	Server             *http.Server
	state              chan *InternalStateChange
}

// // FUNCTIONS

func (a *APIServer) Init(c chan *DBOperationSet, s chan *InternalStateChange, primary DBClient, secondary DBClient) {
	a.Operations = c
	a.state = s
	a.primary = primary
	a.secondary = secondary
	a.lookupOperationSet = map[string]*DBOperationSet{}
	a.Server = &http.Server{
		Addr:    a.ListenAddr,
		Handler: nil,
	}
	a.SetupRoutes()
	if err := a.Server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func (a *APIServer) SetupRoutes() {
	healthHandler := http.HandlerFunc(a.HandleHealth)
	seedDataHandler := http.HandlerFunc(a.HandleSeedData)
	readDataHandler := http.HandlerFunc(a.HandleReadData)
	writeDataHandler := http.HandlerFunc(a.HandleWriteData)
	deleteDataHandler := http.HandlerFunc(a.HandleDeleteData)
	configureClientHandler := http.HandlerFunc(a.HandleConfigureClient)
	reconnectClientHandler := http.HandlerFunc(a.HandleReconnectClient)
	connectClientHandler := http.HandlerFunc(a.HandleConnectClient)
	disconnectClientHandler := http.HandlerFunc(a.HandleDisconnectClient)
	operationHandler := http.HandlerFunc(a.HandleOperation)
	workflowHandler := http.HandlerFunc(a.HandleWorkflow)
	closeHandler := http.HandlerFunc(a.HandleClose)
	queryHandler := http.HandlerFunc(a.HandleQuery)
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "health"), Middleware(healthHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "data/read"), Middleware(readDataHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "data/write"), Middleware(writeDataHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "data/delete"), Middleware(deleteDataHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "data/seed"), Middleware(seedDataHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/configure/primary"), Middleware(configureClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/configure/secondary"), Middleware(configureClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/reconnect/primary"), Middleware(reconnectClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/reconnect/secondary"), Middleware(reconnectClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/connect/primary"), Middleware(connectClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/connect/secondary"), Middleware(connectClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/disconnect/primary"), Middleware(disconnectClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "client/disconnect/secondary"), Middleware(disconnectClientHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "operation/"), Middleware(operationHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "workflow"), Middleware(workflowHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "close"), Middleware(closeHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "query"), Middleware(queryHandler))
}

func (a *APIServer) AddOperationSet(dbo *DBOperationSet) {
	a.lookupOperationSet[dbo.ID] = dbo
	a.Operations <- dbo
}

func (a *APIServer) clearLookupOperationSet() {
	a.lookupOperationSet = make(map[string]*DBOperationSet)
}

func (a *APIServer) LookupOperationSet(guid string) *DBOperationSet {

	if val, ok := a.lookupOperationSet[guid]; ok {
		return val
	}

	return nil
}

// CUSTOM MIDDLEWARE

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(w, r)
	})
}

// HANDLERS - BASE

func (a *APIServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("ok."))
		if err != nil {
			fmt.Println(err)
		}

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleClose(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleClose().")
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusMethodNotAllowed)
	case http.MethodPost:
		go func() {
			time.Sleep(3 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := a.Server.Shutdown(ctx); err != nil {
				a.Server.Close()
			}

			close(a.Operations)
			a.state <- newInternalState("process_then_stop")

			time.Sleep(10 * time.Second)
			a.state <- newInternalState("processing_complete")
		}()

		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleQuery().")
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusMethodNotAllowed)
	case http.MethodPost:
		if !verifyContentType(r, "application/json") {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		q := &DBQueryDataSet{}

		errString := getRequestBody(r, q)
		if len(errString) > 0 {
			switch errString {
			case "default":
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			case "Request body too large":
				http.Error(w, errString, http.StatusRequestEntityTooLarge)
			default:
				http.Error(w, errString, http.StatusBadRequest)
			}
		}
		op := Query(a.primary, q)
		ops := newDBOperationSet(nil)
		ops.AddOperation(op)
		a.AddOperationSet(ops)

		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleWorkflow(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleOperation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleOperation().")

	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", "operation"))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	opID := urlPathSegments[len(urlPathSegments)-1]

	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(a.findOperation(opID))
		if err != nil {
			fmt.Println(err)
		}

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HANDLERS - DATA

func (a *APIServer) HandleSeedData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleSeedData()")
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusMethodNotAllowed)
	case http.MethodPost:
		if !verifyContentType(r, "application/json") {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		ds := &DBDataSet{}
		errString := getRequestBody(r, ds)
		if len(errString) > 0 {
			switch errString {
			case "default":
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			case "Request body too large":
				http.Error(w, errString, http.StatusRequestEntityTooLarge)
			default:
				http.Error(w, errString, http.StatusBadRequest)
			}
		}
		op := SeedData(a.primary, ds)
		ops := newDBOperationSet(nil)
		ops.Operations = append(ops.Operations, op)
		a.Operations <- ops

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleWriteData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}
	case http.MethodPost:
		// do stuff.

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleReadData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}
	case http.MethodPost:
		// do stuff.

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleDeleteData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}
	case http.MethodPost:
		// do stuff.

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HANDLERS - CLIENT

func (a *APIServer) HandleConfigureClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}
	case http.MethodPost:
		// do stuff.

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleConnectClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}
	case http.MethodPost:
		// do stuff.

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleDisconnectClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}
	case http.MethodPost:
		// do stuff.

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *APIServer) HandleReconnectClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTest().")
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(ToJSON(a))
		if err != nil {
			fmt.Println(err)
		}
	case http.MethodPost:
		// do stuff.

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// INTERNAL COMMON HTTP

func interpretHttpError(err error) string {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Sprintf("Request body contains badly-formed JSON")
	case errors.As(err, &unmarshalTypeError):
		return fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Sprintf("Request body contains unknown field %s", fieldName)
	case errors.Is(err, io.EOF):
		return "Request body must not be empty"
	case err.Error() == "http: request body too large":
		return "Request body too large"
	default:
		return "default"
	}
}

func getRequestBody(r *http.Request, output interface{}) string {

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&output)
	if err != nil {
		return interpretHttpError(err)
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return "Request body must only contain a single JSON object"
	}

	return ""
}

func verifyContentType(r *http.Request, ct string) bool {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != ct {
			return false
		}
	}
	return true
}

func (a *APIServer) findOperation(ID string) []byte {

	for i, opset := range a.lookupOperationSet {
		if i == ID {
			return opset.GetOutputJSON()
		}

		for oi, op := range opset.lookupOperation {
			if oi == ID {
				return op.GetOutputJSON()
			}
		}
	}

	return nil
}
