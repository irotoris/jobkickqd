# jobkickqd

Jobkickqd is a job management tool that triggers commands in a job queue.
You can execute a command via Cloud PubSub.

## Requirements

- Google Cloud PubSub
  - 2 topics for a job queue and log queue.
- Google Cloud Credentials
  - create subscription and subscribe topic role and a service account key.
  - e.g.) `export GOOGLE_APPLICATION_CREDENTIALS=<credential_key_path>`
  
## Install

require golang `1.11 or later`

```bash
make
```

## Usage

### daemon(job executor)

```bash
$ jobkickqd daemon \
    --app appName \
    --workDir workDir \
    --jobQueueTopic jobQueueTopic \
    --logTopic logTopic \
    --projectID projectID
```

- `--app`: Daemon executes a command in pubsub message when `--app` name is match. `--app` is unique in all daemon.
- `--workDir`: create a directory in this work directory before command executes

### client(push a command)

```bash
$ jobkickqd submit \
    --app appName \
    --jobQueueTopic jobQueueTopic \
    --logTopic logTopic \
    --projectID projectID \
    --command command \
    --environment "key1=value1,key2=value2,..." \
    --timeout timeout \
    --jobID jobID
```

- `--jobID`: This is a unique id in all job history.

e.g.) exec `uname` command and get command output.

```bash
$ jobkickqd submit --app app1 \
    --jobTopicName jobq \
    --logTopic logs \
    --projectID my_project \
    --command "uname"
    --jobID testjob

INFO[0006] Published a message to pubsub[jobq] with a message ID: 469182461657193
INFO[0007] Job stdout/stderr:
Linux
```

## Build

build a binary, require golang `1.11 or later`

```bash
make build
```

or doker build

```bash
docker build . -t jobkickqd
```