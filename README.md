# Moul

> The minimalist photo collection generator.

Demo: https://demo.moul.app

## Features

- **Simplicity** simple and easy to use
- **Smart** deterministic image layout
- **JAMStack** Future proof [JAMStack](https://jamstack.org)
- **The One** One Binary to rule them all

## Installation

### Advanced

You can download the binary on [release](https://github.com/sophearak/moul/releases) page. Put it in your local bin folder. Make sure your bin folder are in your PATH.

### Detailed

<details><summary>macOS</summary>
<p>

```bash
# download binary
$ wget https://github.com/sophearak/moul/releases/download/v1.0.0-beta.0/moul-darwin

# move binary into local bin
# if ~/.local/bin doesn't exist. $ mkdir -p ~/.local/bin
$ mv moul-darwin ~/.local/bin/moul

# make sure that ~/.local/bin in your $PATH
# if you're using bash
$ echo "export PATH=$PATH:~/.local/bin" >> ~/.bashrc

# if you're using zsh
$ echo "export PATH=$PATH:~/.local/bin" >> ~/.zshrc

# If you don't know what it is. Use bash.

# reload
$ source ~/.bashrc
```

</p>
</details>

<details><summary>Linux</summary>
<p>
I'm 300% sure, you know what to do. See Advanced
</p>
</details>

<details><summary>Windows</summary>
<p>

- **Step 1**: Download the binary [here](https://github.com/sophearak/moul/releases)
- **Step 2**: Create a folder in `C:\bin` and put the downloaded file in there
- **Step 3**: Add `C:\bin` to your `Environment Variables` by
  - **Step 3.1**: right-click `My Computer` -> click `Properties`
  - **Step 3.2**: In the `System Properties` window, click the `Advanced` tab, and then click `Environment Variables`.
  - **Step 3.3**: In the `System Variables` window, highlight `Path`, and click `Edit`.
  - **Step 3.4**: In the Edit `System Variables` window, insert the cursor at the end of the `Variable` value field.
  - **Step 3.5**: If the last character is not a semi-colon (;), add one.
  - **Step 3.6**: After the final semi-colon, add `path C:\bin` -> click `OK`

</p>
</details>

## Usage

```bash
# create new photo collection
$ moul new my-collection

# Place photos into its desire folders

# generate your photo collection
$ cd my-collection && moul dev

# build for production
$ moul build
```

> It depends on how many photos you added to `photos/collection`, the commend `moul dev` might take a while. It's a good time to grab coffee.

## Recommended size

- `photos/cover`: 2560px width - any landscape aspect ratio will work fine
- `photos/profile`: 1024px width - square (1:1 aspect ratio)
- `photos/collection`: 2048px width - any aspect ratio

## Configuration

```json
{
  "site": {
    "url": "https://demo.moul.app",
    "name": "Moul",
    "bio": "The minimalist photo collection generator"
  },
  "social": {
    "twitter": "mouldotco",
    "youtube": "",
    "facebook": "",
    "instagram": ""
  }
}
```

> The link for twitter will be come `https://twitter.com/mouldotco`

## Deployment

You can pretty much deploy the `dist` folder to any static site hosting. That includes

- Firebase Hosting - https://firebase.google.com/docs/hosting
- Netlify - https://netlify.com
- Now - https://zeit.co/now
- Surge - https://surge.sh
- Github page - https://pages.github.com

and more.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

- [github.com/denisbrodbeck/sqip](https://github.com/denisbrodbeck/sqip)
- [github.com/fogleman/primitive](https://github.com/fogleman/primitive)
- [github.com/gobuffalo/packr](https://github.com/gobuffalo/packr)
- [github.com/gobuffalo/plush](https://github.com/gobuffalo/plush)
- [github.com/spf13/cobra](https://github.com/spf13/cobra)
- [github.com/spf13/viper](https://github.com/spf13/viper)
- [github.com/tdewolff/minify](https://github.com/tdewolff/minify)
- [github.com/h2non/bimg](https://github.com/h2non/bimg)
