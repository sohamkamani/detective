# Detective 🔎

Detective is a distributed application health monitoring library. It allows you to monitor arbitrary dependencies in your application, and compose other detective instances to create a distributed monitoring framework.

A typical service oriented architecture looks like this:

![service oriented architecture](images/webapp-arch.png)

You can replace the components with arbirtrary components of your own, but you get the idea. Detective allows you to enable each application to monitor its own dependencies, including dependencies with contain another detective instance. By doing so, you can monitor your infrastructure in a distributed manner, where each service _only_ monitors _it's own_ dependencies.

![service oriented architecture with detective](images/detective-arch.png)


