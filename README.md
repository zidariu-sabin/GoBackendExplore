# GoBackendExplore

### Scope
- Golang project resembling the backend component of a movie reviewing web application, created for the purpose of practicing common industry used concepts such as
	- data seeding
	- database manipulation
		- postgres database
		- migrations
		- transactions
		- crons
			- checking and deleting user tokens if they are expired
	- error logging
	- authentification,
		- statefull auth
	- testing,
	- ci/cd pipelines,
	- rate limiting 
	- image pulling from a cloud store 

### Improvements:
- add invalid payload request errors
- refactor router to use subrouter to avoid code duplication on require user
