/*
Detective is a distributed application health monitoring library.

Detective allows you to simultaneously, and independently monitor the health of any application along with its dependencies.

Using Detective

A Detective instance represents a group of dependencies. A dependency can be anything that the application depends on, and which can return some sort of error. For example, a database whose connection time out, or an HTTP web service that doesn't respond.

	// Create a new detective instance, and add a dependency
	d := detective.New("application")
	dep := d.Dependency("database")

	// Register a function to detect a fault in the dependency
	dep.Detect(func() error {
		// `db` can be an instance of sql.DB
		err := db.Ping()
		return err
	})

	// Create a ping endpoint which checks the health of all dependencies
	http.ListenAndServe(":8080", http.HandlerFunc(d.Handler()))

The ping endpoint will check the health of all dependencies by calling their detector function registered in the "Detect" method

Composing instances

The endpoint that was defined in the previous example can itself be registered in other detective instances:

	d2 := detective.New("dependent_application")

	// The ping handler defined in the last example can be registered as an endpoint
	d2.Endpoint("http://localhost:8080/")

	// The ping handler of `d2` will now call `d`s ping handler and as a result, monitor `d`s dependencies as well
	http.ListenAndServe(":8081", http.HandlerFunc(d2.Handler()))

*/
package detective
