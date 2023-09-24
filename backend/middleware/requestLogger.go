package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"api-lisa/utils/log"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLoggerActivity func for logging
func RequestLoggerActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		writeLogReq(c)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		dur := time.Since(t)
		c.Set("Latency", dur.String())
		if strings.Contains(c.FullPath(), "download") {
			writeLogResp(c, "ok")
		} else {
			writeLogResp(c, blw.body.String())
		}
	}
}

func writeLogReq(c *gin.Context) {
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

		re := regexp.MustCompile(`\r?\n`)
		var request = re.ReplaceAllString(readBody(rdr1), "")
		log.Log.WithFields(logrus.Fields{
			"logType":     "Request",
			"url":         c.Request.URL.Path,
			"method":      c.Request.Method,
			"requestId":   requestid.Get(c),
			"userAgent":   c.Request.UserAgent(),
			"requestBody": request,
		}).Info()
		c.Request.Body = rdr2
	} else {
		if c.FullPath() != "/" {
			log.Log.WithFields(logrus.Fields{
				"logType":   "Request",
				"url":       c.Request.URL.Path,
				"method":    c.Request.Method,
				"requestId": requestid.Get(c),
				"userAgent": c.Request.UserAgent(),
			}).Info()
		}
	}
}

func writeLogResp(c *gin.Context, resp string) {
	latency, _ := c.Get("Latency")
	if c.FullPath() != "/" {
		log.Log.WithFields(logrus.Fields{
			"logType":      "Response",
			"url":          c.Request.URL.Path,
			"method":       c.Request.Method,
			"requestId":    requestid.Get(c),
			"userAgent":    c.Request.UserAgent(),
			"latency":      latency.(string),
			"responseBody": resp,
		}).Info()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

func SetupLogger() {
	log.Log.Println("Setup Logger Start")
	lumberjackLogger := &lumberjack.Logger{
		// Log file absolute path, os agnostic
		Filename:   filepath.ToSlash("./log/log"),
		MaxSize:    5, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}
	multiWriter := io.MultiWriter(lumberjackLogger, os.Stderr)
	Formatter := new(logrus.JSONFormatter)
	// You can change the Timestamp format. But you have to use the same date and time
	Formatter.TimestampFormat = "02-01-2006 15:04:05"

	log.Log.SetFormatter(Formatter)
	log.Log.SetOutput(multiWriter)
	log.Log.Println("Setip Logger Finish")
}
