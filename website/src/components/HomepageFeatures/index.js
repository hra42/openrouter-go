import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Zero dependencies',
    description: (
      <>
        Pure Go standard library. No transitive supply chain, no version
        churn — <code>go get</code> and you're done.
      </>
    ),
  },
  {
    title: 'Full API coverage',
    description: (
      <>
        Chat, legacy completions, streaming, tool calling, structured outputs,
        multimodal inputs (image / audio / PDF / text), embeddings, and the
        Anthropic-compatible Messages endpoint.
      </>
    ),
  },
  {
    title: 'Streaming done right',
    description: (
      <>
        SSE streams with channel-based iteration, context cancellation, and
        automatic retry with exponential backoff.
      </>
    ),
  },
  {
    title: 'Functional options',
    description: (
      <>
        One pattern, everywhere. <code>WithModel</code>,{' '}
        <code>WithTools</code>, <code>WithResponseFormat</code> — discoverable,
        composable, extensible.
      </>
    ),
  },
  {
    title: 'Typed errors',
    description: (
      <>
        <code>*RequestError</code> with <code>IsRateLimitError</code>,{' '}
        <code>IsAuthenticationError</code>,{' '}
        <code>IsContextLengthError</code>, and more.
      </>
    ),
  },
  {
    title: 'Agent-friendly',
    description: (
      <>
        <code>llms.txt</code>, <code>AGENTS.md</code>, and a machine-readable{' '}
        <code>api-surface.json</code> so Claude Code and friends can build on
        the SDK without crawling.
      </>
    ),
  },
];

function Feature({ title, description }) {
  return (
    <div className={clsx('col col--4', styles.feature)}>
      <div className="padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
