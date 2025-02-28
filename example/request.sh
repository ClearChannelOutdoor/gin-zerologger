#/!/bin/sh

# should be a log.Debug
echo should log at debug level
curl -XGET localhost:8080

# should be a log.Info
echo
echo should log at info level
curl -XGET localhost:8080/300

# should be a log.Warn
echo
echo should log at warn level
curl -XGET localhost:8080/400

# should be a log.Error
echo
echo should log at error level and include additional error details
curl -XGET localhost:8080/500

# should be a log.Debug and include request body
echo
echo should log at debug level and include the request body
curl -XPOST localhost:8080/logbody --data "{\"test\":true}"

# should be a log.Debug with correlation ID
echo
echo should log at debug level with correlation ID
curl -XGET -H"X-Correlation-ID:12345" localhost:8080

# should be a log.Debug with context settings
echo
echo should log at debug level with additional gin context values
curl -XGET localhost:8080/logcontext
