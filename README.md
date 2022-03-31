## GOCLI

##### CURRENTLY UNDER DEVELOPMENT

- [Idea](#idea)
- [Contributing](#contributing)

  
<br />

#### Idea

I see Go as being fast, great with concurrency, and has great semantics

This is intentded to be a project that explores the idea of creating a command-line-interface tool in go that can do the following things:
  - Retreive Data from various sources
  - Write Data to files
  - Manipulate Data from files
  - Store Sources, 
  - Store Sequences of Operation
  - View the stored Sources, Sequences, and Data
  - Send Data to destinations
  - Runtime has state that is persistent
  - You can update key/value pairs on state


I want this tool to follow consistent, safe patterns, and it will be optimized for speed later if necessary.

Additionally I want to reduce the amount of maintenence on this project to an absolute minimum. Could this potentially be a truly complete project?

The way I think about building a cli is how easy it is to navigate through submenus, and selecting the correct option.

So for example I'd like to do something like:

user@host:$ gocli

1) pipelines
2) sources
3) sequences
4) view
5) send

--> 2

then this would show me a list of sources that I can choose from
choose a source
then this would show me a list of actions I can perform with the value from that source

The way I think this could be useful if you could reduce the barier to connect to multiple types of storage platforms

Ways you can get data:

1. From API endpoint (http)
2. Mongodb URI
3. SQL URI
4. Copying files scp
5. UDP streaming

Being able to quickly, easily, and consistently reteive, store, and optionally manipulate data.

Also this leads into being able to connect different types of storage mediums by a tool that can translate between them.

This will be my first Go endeavor so if you taking a look at all the code, just know it's a work in progress.

I will try to update this README.md as best I can as things change and progress.

<br />  
#### Contributing
If you would like to contribute to this project email: matt@composersoftware.com

<br />  
