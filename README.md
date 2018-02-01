# IMT2681 - Cloud Technologies - Assignment 1

## Service Specification:

The service has to be deployed on either Google Compute Engine or Heroku and expose an API that commits to the following specifications. The service has to be written in Go programming language, must pass Lint and Vet without warnings, and must have at least 20% test coverage. The service is stateless, should not store or record any information, and it should allow concurrent access from multiple clients at the same time. 

#### Invocation

\- *Base path:* /projectinfo/v1/[url]

\- *Method:* GET

\- *Example:* http://localhost:8080/projectinfo/v1/github.com/apache/kafka

(with apache referring to the organisation or username and kafka as a repository)

#### Response

\- *Response payload:*

{

​    "project": {
​        "type": "string"
​    },
​    "owner": {
​        "type": "string"
​    },
​    "committer": {
​        "type": "string"
​    },

​    "commits": {
​        "type": "number"
​    },
​    "language": {
​        "array": {
​            "items": {
​                "type": "string"
​            }
​        }
​    }
}

#### Example: 

{

​    "project": "kafka",

​    "owner": "apache",

​    "committer": "enothereska",

​    "commits": 19,

​    "language": ["Java", "Scala", "Python", "Shell", "Batchfile"]

}
