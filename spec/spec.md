

SPringMuD是一种新的SPMD架构，追求硬件的极简和高效。概念上是16个小核心并发执行同样的Program，每个小核心被称为一个Yarn（纱线），4个Yarn是一个Cord（绳子），4个Cord是一个Braid（辫子）。每个Core里面有两个Braid，同一时刻只能有一个Braid处于活动状态，外部的控制核可以切换活动的Braid。之所以要有两个Braid，是为了能让计算和DMA能interleave起来。

一个Braid里面的4个Cord是时分复用的，Cord执行指令的频率是系统时钟的1/4。一个Cord内部的4个Yarn总是同时执行的。4个Yarn可以同时发出load/store指令，Data Ram被分为8个Bank，期待这4个Yarn的目标地址分布在不同的Bank里，以避免冲突，如果有冲突，用塞空泡+重启流水线的方法来增加额外的周期。

当前的第一版不支持Cache，不支持多周期指令（div sqrt rem）。

每个Yarn内部的软件可见寄存器包括：32个32位寄存器（其中有4个寄存器是临时寄存器），32个64位寄存器（其中有4个寄存器是临时寄存器），一个Bool寄存器，一个forward branch target寄存器栈，一个YarnID寄存器（只读）。8个临时寄存器采用Flipflop实现，其他56个寄存器均采用SRAM实现。在切换Braid的时候，临时寄存器会被破坏掉（取值不定），其他寄存器的值会被保持。

之所以要引入临时寄存器，是为了：一、减少读写寄存器堆的功耗；二、弥补SRAM所实现的寄存器堆读写口有限的缺点。

每个Braid内有一个PC寄存器，一个返回地址栈以及若干控制寄存器、状态寄存器，被所有的Yarn所共享。

每个Yarn每拍最多可以译码和执行4条指令。采用VLIW的指令编码，每个指令组包含若干指令，每个指令为32位，在指令的编码中，用最高位显式标记指令组的边界。

Branch指令分forward和backward两种。forward branch的实现方法：标记Yarn内部的forward branch target寄存器，如果PC没有达到这个值，则指令被视为Nop。backward branch目前只支持按第0个Yarn的Bool寄存器来判断是否Branch。

不支持跳转到寄存器的取值（indirect branch）。

宏观上采用四级流水的结构：取值&译码（Fetch），访问寄存器堆（AccRF），执行和访问Memory		q（ExMem），以及写回（Write）。

Load指令在写回时，先写回到一个WriteBuf中，等到下一条Load指令的AccRF级，这个WriteBuf才会真正写入到寄存器堆的DRAM中。这样，Load指令的写回和Store指令的读源数据，都统一在AccRF级进行。

一个逻辑上的流水级在实现中包含4个Phase。Fetch级的4个Phase分别完成：一、计算NextPC；二、取指令；三、取指令（因为指令大概率不是4Word对齐的，需要两个周期才能取出来所需的指令）；四、译码。

对于32位的寄存器堆，AccRF的4个Phase分别完成：一、写回32位执行部件的计算结果；二、写回Load指令的WriteBuf或者读取Store指令的源数据；三、读取RB的值或者LoadStore指令的基址寄存器；四、读取RA的值或者LoadStore指令的基址寄存器。

对于64位的寄存器堆，AccRF的4个Phase分别完成：一、写回64位执行部件的计算结果；二、写回Load指令的WriteBuf或者读取Store指令的源数据，或者读取乘累加指令的RA的值；三、读取RB的值；四、读取乘累加指令的RC的值，或者两操作数计算指令的RA的值。

除了MAC指令的latency为两拍之外，其他指令的latency都是一拍。软件来避免两拍的MAC指令所引入的hazard。

具体实现中，Instruction SRAM的宽度不到32位，一些固定为0的域，就被忽略掉了。

一条指令会占用的资源包括：寄存器堆的读口；寄存器堆的写口；Bool寄存器的写口；执行部件（目前四个部件：32位执行部件负责32位整点，64位执行部件负责64位整点以及单精度、双精度浮点，Mem部件负责Load&Store指令，Branch部件负责跳转）。一个指令组内部的指令不能存在资源冲突，这一点要靠编译器或者汇编程序员来保证。

指令编码的格式如下：

| 类型           | OP               | RD       | RA       | RB     |
| -------------- | ---------------- | -------- | -------- | ------ |
| Calc-Reg-Reg   | 9bit             | RD       | RA       | RB     |
| Calc-Reg       | 9bit             | RD       | extra op | RB     |
| Calc-Reg-Imm8  | 6bit op + imm7~5 | RD       | imm4~0   | RB     |
| Load           | 9bit             | RD       | imm4~0   | RB     |
| Store          | 9bit             | RM       | imm4~0   | RB     |
| Bool-Set-on-XX | 9bit             | imm9~5   | imm4~0   | RB     |
| ControlFlow    | 9bit             | imm14~10 | imm4~0   | imm9~5 |
|                |                  |          |          |        |


