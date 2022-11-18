## What are the different _Fields_ used for?
- **name_pattern** | This field must match the blob path from your release. If your blob path includes information about the current version you should replace it with ``((version))``
- **cur_version** | This field is here to not require the version number as part of the blob path, is a bad practice that makes bumping versions dangerous and difficult
- **constraints** | Specify a comma-separated list of constraint as in other package managers. This is useful if you want to embed multiple versions of the same dependency as separate blobs.
- **vers_regexp** | The regular expression that will be used to detect new versions. It will be matched against the content of the webpage specified in `vers_url` field
- **vers_url**    | The URL to a webpage or unauthenticated API where we can get a plain-text doc containing references to the latest available versions of this dependency
- **blob_url**    | The URL for downloading the blob. The text `((version))` will be replaced by the corresponding value at runtime.
- **src_url**     | (Optional) The URL for downloading the source code of the blob. Useful if the blob is a precompiled binary for which we want to keep a copy of the source code
