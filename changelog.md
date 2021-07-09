# Changelogs

## Doctor v0.3.0 (Unreleased)

### Breaking changes

- [ ] The Doctor API is totally revamped and should now be easily usable by other projects. The CLI tool is now built upon this API.
- [x] **Sections** are now called **chapters**. This terminology scales better with larger projects, while still being applicable for smaller ones. The source directory default is now `chapters`, so you'll need to rename your `secs` directory. An alternative is to use the new option to configure the source directory, namely add `source-directory = "secs"` into your `doctor.toml` file.

### New features

- [x] New configuration option: You can now choose the source directory (the directory storing your Markdown files), by using the `source-directory` option under the `[build]` header in `doctor.toml`. By setting `source-directory = "."`, you can thus have all your chapters stored in the document root, if that is desireable. The default value is `"chapters"`.


## Doctor v0.2.7 (2020-07-06)

### New features

- [x] New command: `doctor stats`. Currently only supports a simple wordcount, which you get by running `doctor stats --wordcount`. Does nothing when run without a flag, at the moment. Expect lots of fun statistics about your document here.
- [x] New command: `doctor edit`. Use this to open a section directly in your preferred editor (respects `$EDITOR`). Uses the same section finding method as the `remove` and `move` commands.

### Improvements

- [x] Lua filters are now written to the Doctor data directory (which respects `XDG_DATA_DIR`). This avoids the need for writing and deleting the filters on every build, while still keeping the project directory clean. If filters are changed in a future Doctor update, they will be rewritten.


## Doctor v0.2.6 (2020-06-04)

### New features

- [x] New configuration option `bibliography-file` under the `[bib]` table. This allows you to use another filename than the default `references.bib`. The BibTeX file must still be located in the `assets` directory

### Improvements

- [x] Further improved title recognition. `doctor list` should benefit from not being as dependent on the filenames.

### Bug fixes

- [x] CGo was previously enabled, which dynamically linked libc. I am building Doctor on Arch Linux, which uses a newer version of libc than most Debian distributions. This made Doctor crash on said distributions. CGo is not necessary, and is now disabled.


## Doctor v0.2.5 (2020-05-27)

### New features

- [x] When running `doctor new`, you can supply the `--default` flag to automatically create a classic report structure of abstract, introduction, theory, method, results, discussion and conclusion.

### Improvements

- [x] `doctor new` now automatically fills in the `title`-field in `doctor.toml`, based on the capitalized name of the directory the document is created in. This is only a convenience feature, but should be less confusing than the previous `title = "TITLE"` placeholder.


## Doctor v0.2.4 (2020-05-23)

### New features

- [x] Doctor now has help texts for the main program and all sub-commands. Find them by running `doctor --help` or `doctor <command> --help`.

### Bug fixes

- [x] Table reference numbers were previously found manually by inserting the table count. In the case of table numbers like "Table 2.1" or "Table A.3", this would output incorrectly. This is now fixed by using `\label` and `\ref` instead.
- [x] `doctor list` did not output any errors. Although they happen rarely, this is now fixed.


## Doctor v0.2.3 (2020-05-16)

I'm fearing Equinox might have abandoned their project and am thus moving to updates hosted on GitHub. [go-github-selfupdate](https://github.com/rhysd/go-github-selfupdate) seems to provide the desired functionality, and this should be properly working now.

You will need to download this version manually and place it in your `PATH`, but when this is done, you can update your binary with `doctor --update`.

If you have Doctor installed with the Go tool, you will follow the HEAD of this repo, and will need to run `go install -u github.com/kmaasrud/doctor` to update.


## Doctor v0.2.2 (2020-04-25)

### New features

- [x] Doctor now comes with a couple of CSL styles embedded, which are written into the `assets` folder when required.
- [x] **IMPORTANT**: Lots of TOML options are added and the structure is changed. Refer to [the docs](https://kmaasrud.com/doctor/config) to see what's new and adapt your documents' `doctor.toml` to the new headers.


## Doctor v0.2.1 (2020-04-21)

### Bug fixes

- [x] Citeproc interfered with the Lua filter for cross-referencing. This is now fixed.
- [x] Table cross-referencing was not working because of a indexing error which is now corrected.


## Doctor v0.2.0 (2020-04-21)

Doctor is now published to Equinox. This means it has multiple installers and supports self-updating. You can replace your previous binary or Go install with the new one on the [download page](https://dl.equinox.io/kmaasrud/doctor/stable), and update to any new version with `doctor --update`.

### New features

- [x] Command: `doctor move` to change the position of sections.
- [x] Command: `doctor list` to list the sections of your document, along with their index.
- [x] Command: `doctor --update` to update Doctor to the newest version.

### Improvements

- [x] Added options to `doctor.toml`:
    - [x] Boolean option: Allow turning off all Lua filters with the `lua-filters`.
    - [x] Boolean option: Specify if you want numbered sections or not with `number-sections`.
    - [x] String option: Specify an inclusion into the LaTeX header with `latex-header`.
        
        > This option only has an effect when the output format is `pdf`. For now that is the only option, so you won't have to worry about this just yet.

    - [x] String option: Specify a title for the section containing your references with the `references-title` option.
- [x] Even though it is mostly unused for now, the Doctor data dir now respects `$XDG_DATA_DIR` and `%APPDATA%`.
- [x] Big revamp to documentation and lots of stuff added to it.

### Bug fixes

- [x] Important error messages should no longer be swallowed by the Doctor messaging system.
- [x] Directory `src` was added as a resource path. This is now changed to `secs`.
- [x] Doctor forced you to have Tectonic installed to properly function. This is no longer the case.


## Doctor v0.1.2 (2020-04-15)

-  Added versioning feature. Run `doctor --version` or `doctor -v` to see which version you are currently on.


## Doctor v0.1.1 (2020-04-14)

Fixes issue with a malformed file `'` being part of the repo. You should now be able to install with

    go install github.com/kmaasrud/doctor@latest

Other than that, the binaries are exactly the same as in `v0.1.0`


## Doctor v0.1.0 (2020-04-13)

KODB is now rebranded as Doctor and written totally from scratch in Go. Say hello to a document builder that is blazingly fast, safe and most importantly a pleasure to use.

- [x] `doctor` command line tool that parses commands supplied to it.
- [x] Command: `doctor new` for creating a new document.
- [x] Command: `doctor build` for building a PDF from the Doctor document (currently only supporting [Tectonic](https://tectonic-typesetting.github.io/))
- [x] Command: `doctor add` for adding sections.
- [x] Command: `doctor remove` for removing sections.
- [x] Smashing new documentation hosted on it's own webpage.
- [x] TOML file for configuring attributes and compile options.
