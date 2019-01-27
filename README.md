# moul •

> The minimalist photo collection generator.

Demo: https://demo.moul.app

## Features

- **Simplicity** simple and easy to use
- **Smart** deterministic image layout
- **JAMStack** Future proof [JAMStack](https://jamstack.org)
- **The One** One Binary to rule them all

## Installation

Download a [single binary](https://github.com/sophearak/moul/releases), add to `$PATH` and you’re done

### Don't know what todo?

<details><summary>macOS</summary>
<p>

```bash
curl -s https://moul.app/install.sh | sh
```

</p>
</details>

<details><summary>Linux</summary>
<p>
coming soon
<!-- ```bash
curl -s https://moul.app/install.sh | sh
``` -->

</p>
</details>

<details><summary>Windows</summary>
<p>
coming soon
</p>
</details>

## Usage

```bash
# Create new photo collection
$ moul new my-collection

# Place photos into its desire folders

# Add your information in config.json

# Generate your photo collection
$ cd my-collection && moul dev

# Build
$ moul build
```

> It depends on how many photos you added to `photos/collection`, the command `moul dev` might take a while. It's a good time to grab coffee.

## Recommended size

- `photos/cover`: Any landscape aspect ratio will work fine (recommended: 2560px wide)
- `photos/profile`: Square - 1:1 aspect ratio (recommended: 1024px wide)
- `photos/collection`: Any aspect ratio (recommended: 2048px wide)

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
