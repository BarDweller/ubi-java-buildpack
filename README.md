# ubi-java-buildpack

A Simple helper buildpack that configures a jvm installed by the ubi-java-extension.

Extensions are unable to manipulate layers within the app image, and paketo configures
jvm's via app layers, so the ubi-java-extension is unable to install and configure the 
jvm as it has no access to create layers. 

Instead, the extension creates two files into the builder image that are detected by
this buildpack, and used to perform the configuration.

The files are:

- `/bpi.paketo.ubi.java.version` - conveys the version of java installed
- `/bpi.paketo.ubi.java.helpers` - conveys the helpers that would have been used to configure java