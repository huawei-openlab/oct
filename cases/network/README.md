# Network-Test Specifications

Application containers on Linux are a rapidly evolving area, and within this space networking is a particularly unsolved problem, as it is highly environment-specific. 
To solve the container network problems, there are a few existing multi-host networking solutions, such as pipework, weave, flannel, etc. 
Along with the user demands, more network solutions will appear in the future.
However, what is the difference among these soulutions, can they satisfy user demands?  
There is no obvious answer, and there is no effective ways to justify them. 
What we want to do is provding a network test case suit and a auto test framework, 
helping users to get a comparative analysis and choose the proper container network solutions to use.

## Test case design

Container network test case will be designed from two dimensions, test methods and container network model. In different container network circumstance, different test methods will be used to mensurate the network performance. For example, we can test weave network performance through iperf.

### Test methods

Container network will be tested from two aspects. Using ping to test the basic connectivity of the container network. Using iperf to test the performance of container network.

### Container network model
There are a few existing multi-host networking solutions. The basic solution is NAT, which works by hiding the the containers behind a Docker Host IP address. 
The main defect of NAT solution is port conflict, to solve the problem, nat can be bypassed by L3 encapsulation protocols, this is overlay solution.
We will test the performance of these network solutions, to get a comparative analysis and get a base line data. When new network solution appears, we can use the base line data to estimate the new  network solution.

## Test tool development
To test the container network, we will develop some performance testing tools. At first, we can choose existing open source tools, such as iperf. Along with the evolution of container network, corresponding test tools will be developed to statisfy the particular circumstances.

## Future
Please refer to [roadmap](ROADMAP.md) for more information.

## Copyright and license
Code and documentation copyright Huawei Technologies Co., Ltd. Code released under the Apache 2.0 license. Docs released under Creative commons.
