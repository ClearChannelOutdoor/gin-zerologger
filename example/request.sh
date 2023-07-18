#/!/bin/sh

# should be a log.Debug
echo should log at debug level
curl -XGET -H"X-Correlation-ID:testing" localhost:8080

# should be a log.Warn
echo
echo should log at warn level
curl -XGET -H"X-Correlation-ID:testing" localhost:8080/400

# should be a log.Error
echo
echo should log at error level and include additional details
curl -XGET -H"X-Correlation-ID:testing" localhost:8080/500

# should be a log.Debug
echo
echo should log at debug level and include the request body
curl -XPOST -H"X-Correlation-ID:testing" localhost:8080/logbody --data "{\"test\":true}"

# should be a log.Debug
echo
echo should log at debug level
curl -XGET -H"X-Correlation-ID:testing" localhost:8080
