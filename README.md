[![CI](https://github.com/deni1688/gie/actions/workflows/ci.yml/badge.svg)](https://github.com/deni1688/gie/actions/workflows/ci.yml)
# gie (git issue extractor)
cli tool to create multiple issues for a git provider

### Overview 
This tool should make it simple to extract issues to a git hosting service like GitHub or GitLab from your code base.
Putting extraction one command away will hopefully encourage developers to submit more commits with targeted issues attached. As a result, 
tacking different types of changes should become more transparent and measureable. 

### How it works

By prefixing comments with a specified tag inside your code and then running `gie` on that file or directory 
the comment(s) are extracted, formatted, and submitted to your git hosting service. The resulting issue(s) links are then
appended to the commented line. Each issue will also have a reference to the file where it was extracted from in your git hosting service.

### Example

##### Before

```go
    // Issue: Make it possible to do XYZ in refactorMe 
    func refactorMe() {
        // ...        
    }
```

##### Run
```bash
gie -path .
```


##### After
```go
    // Issue: Make it possible to do XYZ refactorMe refactorMe -> closes https://github.com/owner/project/issues/12
    func refactorMe() {
        // ...        
    }
```

#### Post issue creation
After the `gie` command is ran it will also check if any webhooks are configured and send all extracted issues to the
configured endpoints.


#### Configuration

The default config file is sourced from your $HOME/.config/gie.json. You can generate this file by running `gie -setup`.
This will create the config json with the following options:

```json
{
    "host": "String<GITHUB|GITLAB>",
    "token": "String<YOUR_AUTH_TOKEN>",
    "prefix": "String<PREFIX e.g. // TODO>",
    "query": "String<QUERY_FOR_YOUR_REPOS>",
    "webhooks": "List<YOUR_WEBHOOKS_URLS>",
    "exclude": "List<DIRS_FILE_TO_EXCLUDE>"
}
```

You can also override the config on each run by providing a path to another json file using the -config flag. 

### Roadmap

1. Create a GUI version that can be triggered via a flag which will allow selecting which issues should be created and providing better descriptions.



