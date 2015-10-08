## Specs Validator       
      
This file introduce the principle of the runtimeValidator    
      

### Composition     
![Composition](docs/static/composition.PNG "composition map")      

### SequenceMap     
![Sequence](docs/static/sequence.PNG "sequence map")   

### How to Contribute cases    

- What is testsuite
  
  opencontainer/specs pkg conclude two top level structs, LinuxSpec and LinuxRuntimeSpec.     
  The runtimeValidator takes the top level stuct of LinuxSpec or LinuxRuntimeSpec as a testing obj,       
  and create a testsuit for it on to one.    
  For example, spec.Version is a testsuite.

- What is testcase   
   
  Each testsuite should have serveral testcases to support test, the amount of the testcases is     
  depend on the composition cases of input set of the testing obj.
  
- Write testcase   
  Anyone who want to know how to write testcase can go throght the [cases/specversion](cases/specversion)
