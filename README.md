## What are the different _Fields_ used for?
- **name_pattern** | This field must match the blob path from your release. If your blob path includes information about the current version use ``((version))`` as a placeholder in this file
- **cur_version** | This field is here to not require the version number as part of the blob path, is a bad practice that makes bumping versions dangerous and difficult
- **constraints** | Specify a comma-separated list of constraint as in other package managers. This is useful if you want to embed multiple versions of the same dependency as separate blobs.
- **vers_regexp** | The regular expression that will be used to detect new versions. It will be matched against the content of the webpage specified in `vers_url` field
- **vers_url**    | The URL to a webpage or unauthenticated API where we can get a plain-text doc containing references to the latest available versions of this dependency
- **blob_url**    | The URL for downloading the blob. The text `((version))` will be replaced by the corresponding value at runtime.
- **src_url**     | (Optional) The URL for downloading the source code of the blob. Useful if the blob is a precompiled binary for which we want to keep a copy of the source code

## Why including version number in blob path is bad practice?

When you upload a release to a BOSH Director your packaging scripts can interact with the blobs.<br/>
Blobs are available to you as normal files having the same name and path as the blob path itself.<br/>
This means that when bumping a blob, you need to be very careful and review your packaging scripts<br/>
as there might be some hardcoded references to your blob path that you need to update.<br/>
The same thing can happen in configuration files, job templats, etc. They can be anywhere.<br/>

## How do I know which version I'm using without adding it to the blob path?

That's the goal of the `cur_version` field listed above. When a blob gets bumped, that field gets updated
but the blob path remains the same, ensuring that packaging scripts and other configurations still work as expected
without sacrificing the transparency of knowing which version is embedded in the release.

That, together with the `blob_url` and the `sha` field in your `config/blobs.yml` gives you perfect traceability of
the blob's origin, validity and current version whenever is needed.

