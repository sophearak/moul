# Moul

> The minimalist photo collection generator.

Demo: https://demo.moul.app

## Features

- **Simplicity** simple and easy to use
- **Smart** deterministic image layout
- **JAMStack** Future proof [JAMStack](https://jamstack.org)
- **The One** One Binary to rule them all

## Installation

You can download the binary [here](https://github.com/sophearak/moul/releases).

## Usage

```bash
# Create new photo collection
moul new my-collection

# Generate your photo collection
moul dev

# Build for production
moul build
```

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

## Recommended size

- `photos/cover`: 2560px width - any landscape aspect ratio will work fine
- `photos/profile`: 1024px width - square (1:1 aspect ratio)
- `photos/collection`: 2048px width - any aspect ratio

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
