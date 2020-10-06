observability day 1
-------------------

Learning topics:
* Observability
* Implementing metrics
* Scaping metrics
* Write an alert
* SRE/SLOs
* Create an SLO

observability
-------------

When running containers in a distributed environment like kubernetes a proper metrics and monitoring system becomes an absolute necessity. With pods constantly coming and going you need some unified method of determining whether they are functioning correctly.

Understanding the state of your infrastructure and systems is essential for ensuring the reliability and stability of your services. Information about the health and performance of your deployments not only helps your team react to issues, it also gives them the security to make changes with confidence. One of the best ways to gain this insight is with a robust monitoring system that gathers metrics, visualizes data, and alerts when things appear to be broken.


* [Gathering metrics and 4 Golden Signals](https://www.digitalocean.com/community/tutorials/gathering-metrics-from-your-infrastructure-and-applications)
* [USE and RED](https://orangematter.solarwinds.com/2017/10/05/monitoring-and-observability-with-use-and-red/)
* [Intro to metrics and monitoring](https://www.digitalocean.com/community/tutorials/an-introduction-to-metrics-monitoring-and-alerting)

deploy applications
-------------------

Make sure the prometheus service account has cluster reader priviledges for the cluster in order to discover the pods running in each namespace:

    oc adm policy add-cluster-role-to-user cluster-reader -z prometheus-k8s -n openshift-monitoring


Create your namespace

    oc new-project my-namespace
    oc edit ns my-namespace

Add label needed for Prometheus Operator resources

    labels:
      openshift.io/cluster-monitoring: "true"


Install the appp

    oc login
    helm upgrade --install cockroachdb helm-charts/cockroachdb -n my-namespace
    helm upgrade --install go-observe-cockroachdb helm-charts/go-observe-cockroachdb -n my-namespace
    
Find the routes in your project namespace and navigate to the metrics pages
* cockroachdb /_status/vars
* go-observe-cockroachdb /metrics

implementing metrics
--------------------

Review observe.go and understand the 4 contained metrics:

* fakeHttpReqs
* gauge
* insertHistogram
* selectHistogram

Understand how each of these metric works by reviewing [prometheus metrics types](https://prometheus.io/docs/concepts/metric_types/)


scraping metrics
----------------

Understand how prometheus scrapes metrics and the different ways you can configure it. See this [getting started](https://github.com/prometheus/prometheus/blob/master/docs/getting_started.md) guide for raw configuration.


Review [prometheus-operator](https://coreos.com/blog/the-prometheus-operator.html) and understand how service monitors will assist us with automating the raw config.

Review the service monitors in the cockroachdb and go-observe-cockroachdb charts.

create an alert
---------------

Now that we have our apps up and running we should be able to create an alert on either cockroachdb or the go-observe-cockroachdb. Use knowledge of the 4 signals. Use one of the PrometheuRule files located in either helm chart. Have it reviewed once you are finished.


google sre
----------

It’s impossible to manage a service correctly, let alone well, without understanding which behaviors really matter for that service and how to measure and evaluate those behaviors. To this end, we would like to define and deliver a given level of service to our users, whether they use an internal API or a public product.

Books

* [Books](https://landing.google.com/sre/books/)
* [SLOs](https://landing.google.com/sre/sre-book/chapters/service-level-objectives/)
* [Monitoring Distributed Systems](https://landing.google.com/sre/sre-book/chapters/monitoring-distributed-systems/) Defines 4 Golden signals here
* [Practical Alerting](https://landing.google.com/sre/sre-book/chapters/practical-alerting/)
* [Alerting on SLOs](https://landing.google.com/sre/workbook/chapters/alerting-on-slos/)

implement an slo
----------------

Using the references above create an SLO. There's already a partially started one in go-observe-cockroachdb. Have it reviewed once you finish.

extra credit
------------

Add a new metric go the go-observe-cockroachdb app and create an alert/SLO on it.