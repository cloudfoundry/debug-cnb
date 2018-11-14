# `debug-buildpack`
The Cloud Foundry Debug Buildpack is a Cloud Native Buildpack V3 that enables the debuging of JVM applications.

## Detection
The detection phase passes if

* `$BP_DEBUG` exists and build plan contains `jvm-application`
  * Contributes `debug` to the build plan

## Build
If the build plan contains

* `debug`
  * Contributes debug configuration to `$JAVA_OPTS`
  * If `$BPL_DEBUG_PORT` is specified, configures the port the debug agent will listen on.  Defaults to `8000`.
  * if `$BPL_DEBUG_SUSPEND` is specified, configures the JVM to suspend execution until a debugger has attached.  Note, you cannot ssh to a container until the container has decided the application is running.  Therefore when enabling this setting you must also push the application using the parameter `-u none` which disables container health checking.  Defaults to `n`.

## Creating SSH Tunnel
After starting an application with debugging enabled, an SSH tunnel must be created to the container.  To create that SSH container, execute the following command:

```bash
$ cf ssh -N -T -L <LOCAL_PORT>:localhost:<REMOTE_PORT> <APPLICATION_NAME>
```

The `REMOTE_PORT` should match the `port` configuration for the application (`8000` by default).  The `LOCAL_PORT` can be any open port on your computer, but typically matches the `REMOTE_PORT` where possible.

Once the SSH tunnel has been created, your IDE should connect to `localhost:<LOCAL_PORT>` for debugging.

![Eclipse Configuration](eclipse.png)

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

