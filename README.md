# OCT: Open Container Testing

The OCT project aims to promote the [Open Container Initiative](http://www.opencontainers.org/) by providing a universal open container testing `lib/tools`.

## OCT scope
Following [the OCI Principles](https://github.com/opencontainers/specs): 
```
Define a unit of software delivery called a Standard Container. 
The goal of a Standard Container is to encapsulate a software component 
and all its dependencies in a format that is self-describing and portable, 
so that any compliant runtime can run it without extra dependencies, 
regardless of the underlying machine and the contents of the container.
```

OCT covers two areas:
- [Standard Container](#standard-container) 
- [Compliant Runtime](#compliant-runtime)

###Standard Container
A standard container should be a [bundle](https://github.com/opencontainers/specs/blob/master/bundle.md) with one standard 'config.json', one standard 'runtime.json' and one standard 'rootfs'.

OCT provides a [Standard Container Validator](tools/stdContainerValidator/README.md) tool to varify if a bundle could be called `a standard container`.

###Compliant Runtime
A compliant runtime should be the one which could run a standard container `correctly`, either runs `directly` or `indirectly`.

`Correctly` means running by a runtime, all the mounts, uid, and other informations should be exactly same with what defined in config.json/runtime.json.

`Directly` means a runtime could run a standard container without any extra action, just like 'runC'.

`Indirectly` means a runtime(runX) could not run a standard container directly. RunX needs to get a runX-bundle converted from an oci-bundle first and then runs runX-bundle. [Converting tools](#converting-tools)

OCT provides a [Compliant Runtime Validator](tools/specsValidator/README.md) tool to varify if a runtime could be called `a compliant runtime`. (OCT tests on 'runC' so far.)

####Converting tools
One implementaion of converting from OCI to ACI is hosted at: [oci2aci](https://github.com/huawei-openlab/oci2aci)

## Getting Started

- Fork the repository on GitHub
- Read the [HowTO](tools/HowTO.md) for build and test instructions
- Play with the project, submit bugs, submit patches!

### How to involve
If any issues are encountered while using the oct project, several avenues are available for support:
<table>
<tr>
	<th align="left">
	Issue Tracker
	</th>
	<td>
	https://github.com/huawei-openlab/oct/issues
	</td>
</tr>
<tr>
	<th align="left">
	Google Groups
	</th>
	<td>
	https://groups.google.com/forum/#!forum/oci-testing
	</td>
</tr>
</table>


## Who should join
- Open Container project developer/user

### Changes
The `engine` part is now moved to [oct-engine](https://github.com/huawei-openlab/oct-engine)
The `cases` part is now moved to [oct-engine/cases](https://github.com/huawei-openlab/oct-engine/cases)
