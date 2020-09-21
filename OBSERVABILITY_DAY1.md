observability day 1
-------------------

Learning topics:
* Observability
* SRE
* SLOs
* Implementing metrics
* Scaping metrics
* Alerting on the 4 golden signals
* Creating an SLOs

observability
----------------------

When running containers in a distributed environment like kubernetes a proper metrics and monitoring system becomes an absolute necessity. With pods constantly coming and going you need some unified method of determining whether they are functioning correctly.

Understanding the state of your infrastructure and systems is essential for ensuring the reliability and stability of your services. Information about the health and performance of your deployments not only helps your team react to issues, it also gives them the security to make changes with confidence. One of the best ways to gain this insight is with a robust monitoring system that gathers metrics, visualizes data, and alerts when things appear to be broken.

* [Intro to metrics and monitoring](https://www.digitalocean.com/community/tutorials/an-introduction-to-metrics-monitoring-and-alerting)
* [Gathering metrics and 4 Golden Signals](https://www.digitalocean.com/community/tutorials/gathering-metrics-from-your-infrastructure-and-applications)


google sre
----------

A best practices guide for SRE can be found here
https://landing.google.com/sre/books/
https://landing.google.com/sre/sre-book/chapters/practical-alerting/

Use these resources if you want to systematically increase the reliability of your departments services. Great way to get a promotion if your team isn’t already doing it!

links
-----

https://promcon.io/2019-munich/slides/practical-capacity-planning-using-prometheus.pdf
https://grafana.com/blog/2018/08/02/the-red-method-how-to-instrument-your-services/