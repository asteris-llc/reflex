# Reflex

Reflex is an event handling system for Mesos. It's designed for stream
processing workloads built from arbitrary commands communicating over stdio.

An example: reading logs and sinking them to a database.

First, we'll register our handler (called a Task in Reflex):

```
$ curl reflex.service.consul:4000/1/tasks -X POST -d '{"ID":"log-sinker","subscribesTo":["applicationLogLine"],"image":"TODO/oursampleimage","cpu":0.25,"mem":512}'
```

Now, whenever an `applicationLogLine` event comes down the pipe the task is
triggered. We can send an event like this:

```
$ curl reflex.service.consul:4000/1/events -X POST -d '{"type":"applicationLogLine","payload":"this is a log line!"}'
```

You can do that and see that reflex has created a task in Mesos for that event.
Hooray!

As a slightly more complex example, you could register several tasks for
different handlers in Reflex, then each task triggers a new event, so that you
end up with a pipeline of streams of data that will scale dynamically to the
load given it.
