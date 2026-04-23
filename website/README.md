# openrouter-go docs site

Docusaurus-powered website for https://github.com/hra42/openrouter-go.

## Local dev

```bash
cd website
npm install
npm run start
```

The site sources docs from `../docs/` — there is no content duplication. Edit
files under `docs/recipes/` and the site picks them up.

## Build

```bash
npm run build
```

Output goes to `website/build/`.

## Deploy

Pushes to `main` that touch `docs/**` or `website/**` trigger the
`.github/workflows/deploy-docs.yml` workflow, which builds the site and
publishes to GitHub Pages.
