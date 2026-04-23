import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import CodeBlock from '@theme/CodeBlock';
import HomepageFeatures from '@site/src/components/HomepageFeatures';

import styles from './index.module.css';

const quickstart = `package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/hra42/openrouter-go"
)

func main() {
    client := openrouter.NewClient(
        openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
    )

    resp, err := client.ChatComplete(context.Background(),
        []openrouter.Message{openrouter.CreateUserMessage("Hello")},
        openrouter.WithModel("openai/gpt-4o-mini"),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(resp.Choices[0].Message.Content)
}
`;

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <h1 className="hero__title">{siteConfig.title}</h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.install}>
          <CodeBlock language="bash">go get github.com/hra42/openrouter-go</CodeBlock>
        </div>
        <div className={styles.buttons}>
          <Link className="button button--primary button--lg" to="/docs/recipes/getting-started">
            Get started
          </Link>
          <Link
            className="button button--secondary button--lg"
            href="https://pkg.go.dev/github.com/hra42/openrouter-go">
            Godoc
          </Link>
          <Link
            className="button button--secondary button--lg"
            href="https://github.com/hra42/openrouter-go">
            GitHub
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home() {
  return (
    <Layout
      title="Zero-dependency Go client for OpenRouter"
      description="Complete Go bindings for the OpenRouter API: chat, streaming, tools, structured outputs, multimodal, embeddings.">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <section className={styles.quickstart}>
          <div className="container">
            <h2>Quick start</h2>
            <CodeBlock language="go">{quickstart}</CodeBlock>
          </div>
        </section>
      </main>
    </Layout>
  );
}
