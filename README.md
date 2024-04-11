# this fork - Ludo:
- i18n
- update 'Game Controller DB' 
- favorites

# ludo ![Build Status](https://github.com/libretro/ludo/workflows/CI/badge.svg) [![GoDoc](https://godoc.org/github.com/libretro/ludo?status.svg)](https://godoc.org/github.com/libretro/ludo)

Ludo is a work in progress libretro frontend written in go.

<img src="https://raw.githubusercontent.com/kivutar/ludo-assets/master/illustration.png" />

It is able to launch most non GL libretro cores.

It works on OSX, Linux, Linux ARM and Windows. You can download releases [here](https://github.com/libretro/ludo/releases)

## Dependencies

- GLFW 3.3
- OpenGL >= 2.1
- OpenAL

#### On OSX

You can execute the following command and follow the instructions about exporting PKG_CONFIG

    brew install openal-soft

#### On Debian or Ubuntu

    sudo apt-get install libopenal-dev xorg-dev golang

#### On Raspbian

You need to enable the experimental VC4 OpenGL support (Full KMS) in raspi-config.

    sudo apt-get install libopenal-dev xorg-dev

#### On Alpine / postmarketOS

    sudo apk add musl-dev gcc openal-soft-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev mesa-dev

#### On Windows

Setup openal headers and dll in mingw-w64 `include` and `lib` folders.

## Building

    git clone --recursive https://github.com/libretro/ludo.git
    cd ludo
    go build

For more detailed build steps, please refer to [our continuous delivery config](https://github.com/libretro/ludo/blob/master/.github/workflows/cd.yml).

## Running

    ./ludo


## Command goi18n

The goi18n command manages message files used by the i18n package.

```
go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
goi18n -help
```

### Extracting messages

Use `goi18n extract` to extract all i18n.Message struct literals in Go source files to a message file for translation.

```toml
# active.en.toml
[PersonCats]
description = "The number of cats a person has"
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
```

### Translating a new language

1. Create an empty message file for the language that you want to add (e.g. `translate.es.toml`).
2. Run `goi18n merge active.en.toml translate.es.toml` to populate `translate.es.toml` with the messages to be translated.

   ```toml
   # translate.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hello {{.Name}}"
   ```

3. After `translate.es.toml` has been translated, rename it to `active.es.toml`.

   ```toml
   # active.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hola {{.Name}}"
   ```

4. Load `active.es.toml` into your bundle.

   ```go
   bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
   bundle.LoadMessageFile("active.es.toml")
   ```

### Translating new messages

If you have added new messages to your program:

1. Run `goi18n extract` to update `active.en.toml` with the new messages.
2. Run `goi18n merge active.*.toml` to generate updated `translate.*.toml` files.
3. Translate all the messages in the `translate.*.toml` files.
4. Run `goi18n merge active.*.toml translate.*.toml` to merge the translated messages into the active message files.

```bash
goi18n extract -outdir ./i18n/ && cd ./i18n/ && goi18n merge active.*.toml
#....translate -> all translate.*.toml
goi18n merge active.*.toml translate.*.toml && cd ../
```


## License

go-i18n is available under the **MIT license**.
