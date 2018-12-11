---
date: 2017-03-01T11:56:36+08:00
title: 和爸爸一起吹口琴
draft: true
---
# The Tao of HashiCorp

The Tao of HashiCorp is the foundation that guides our vision, roadmap, and product design. As you evaluate using or contributing to HashiCorp's products, it may be valuable to understand the motivations and intentions for our work.

## Workflows, not Technologies

The HashiCorp approach is to focus on the end goal and workflow, rather than the underlying technologies. Software and hardware will evolve and improve, and it is our goal to make adoption of new tooling simple, while still providing the most streamlined user experience possible. Product design starts with an envisioned workflow to achieve a set goal. We then identify existing tools that simplify the workflow. If a sufficient tool does not exist, we step in to build it. This leads to a fundamentally technology-agnostic view — we will use the best technology available to solve the problem. As technologies evolve and better tooling emerges, the ideal workflow is just updated to leverage those technologies. Technologies change, end goals stay the same.

## Simple, Modular, Composable

The Unix philosophy is widely known for preaching the virtues of software that is simple, modular and composable. This approach prefers many smaller components with well defined scopes that can be used together. The alternative approach is monolithic, in which a single tool has a nebulous scope that expands to encompass new features and capabilities. We like to think of the components as blocks that are functional on their own, and can be combined in new and innovative ways.

The simple, modular, composable approach allows us to build products at a higher level of abstraction. Rather than solving the holistic problem, we break it down into constituent parts, and solve those. We build the best possible solution for the scope of each problem, and then combine the blocks to form a solid, full solution.

## Communicating Sequential Processes

Given our belief in simple, modular, composable software, we have several tenets for combining those pieces into a connected system. Communicating Sequential Processes (CSP) is a model of computation wherein autonomous processes connected via a network are able to communicate. We believe that the CSP approach is necessary for managing complexity and building robust scalable systems in a Service Oriented Architecture. Each service should be treated as an individual process that then communicates with other services via an API.

CSP represents the best known way to write software and organize services together to form a system. Nature itself provides the best examples; even the human body is a system of interconnected services — respiratory, cardiovascular, nervous, immune, etc.

## Immutability

Immutability is the inability to be changed. This is a concept that can apply at many levels. The most familiar implementation of immutability is version control systems; once code is committed, that commit is forever fixed. Version control systems, such as git, enjoy widespread use because they offer tremendous benefits. Code becomes versionable, allowing rollback and roll forwards. You can inspect and write code atomically. Using versions enables auditing and creates a clear history of how the current state was reached. When something breaks, the origin of the error can be determined using the version history.

The concept of immutability can be extended to many aspects of infrastructure — application source, application versions, and server state. We believe that using immutable infrastructure leads to more robust systems that are simpler to operate, debug, version and visualize.

## Versioning through Codification

Codification is the belief that all processes should be written as code, stored, and versioned. Operations teams have historically relied on oral tradition to pass along the knowledge of how to build, upgrade and triage infrastructure. But information was easily lost or hidden from the people who needed it most. Codification of critical knowledge promotes information sharing and prevents data loss, as any changes to process are automatically stored and versioned.

HashiCorp products are all designed to follow the codification of knowledge paradigm. Any changes a user makes are versioned and stored to keep a clean history of process.

## Automation through Codification

System administration typically requires an operator to manually make changes to infrastructure, making the position's responsibilities difficult to scale. The scale of infrastructure under management is forever increasing, and manual system administration techniques have struggled to match this new scale. Automation to manage more systems with less overhead is the only option.

While there are many approaches to automation, we promote codification. Codification allows for knowledge to be executed by machines, but still readable by operators. Automated tooling allows operators to increase their productivity, move quicker, and reduce human error. Machines can automatically detect, triage and resolve issues.

## Resilient systems

Resilient systems are built to withstand unexpected inputs and outputs. To accomplish this, the system must have a desired state, a method to collect information on the current state, and a mechanism to automatically adjust the current state to return to the desired state.

We believe that applying this sort of systems rigor to infrastructure is critical to achieve the highest levels of reliability. HashiCorp products will always recognize a desired state through codified knowledge. They will collect real-time information through functionally independent components. And they will provide the tooling to self-heal and auto-recover.

## Pragmatism

We strongly believe in the value of pragmatism when approaching any problem. Many of the principles we believe in like immutability, codification, automation, and CSP are ideals which we aspire towards and not requirements that we dogmatically impose. There are many situations in which the practical solution requires reevaluating our ideals.

Pragmatism allows us to assess new ideas, approaches, and technologies and how they may be adopted to improve HashiCorp's best practices. It would be a mistake to view this as a compromise of first principles, but rather open-mindedness and humility to accept that we may be wrong. The ability to adapt is critical to innovation and one we take pride in.









# HashiCorp的道
2014年12月8日 MITCHELL HASHIMOTO
HashiCorp的道是指导我们的愿景，路线图和产品设计的基础。 当您评估使用或贡献于HashiCorp的产品时，了解我们工作的动机和意图可能很有价值。

## 工作流程，而不是技术

HashiCorp的方法是专注于最终目标和工作流程，而不是底层技术。 软件和硬件将不断发展和完善，我们的目标是使新的工具简单化，同时仍然提供最流畅的用户体验。 产品设计从设想的工作流程开始，以实现设定的目标。 然后我们识别现有的简化工作流的工具。 如果一个足够的工具不存在，我们将进行构建。 这导致了与技术无关的观点 - 我们将使用最好的技术来解决问题。 随着技术的发展和更好的工具出现，理想的工作流程只是更新，以利用这些技术。 技术变革，终极目标保持不变。

## 简单，模块化，可组合

Unix的理念广为人知，讲述了简单，模块化和可组合的软件的优点。 这种方法喜欢许多较小的组件，具有可以一起使用的明确定义的范围。 替代方法是整体的，其中单个工具具有不均匀的范围，扩展到包含新的特征和功能。 我们喜欢将组件视为自己功能的组件，并且可以以新的和创新的方式组合。

简单，模块化，可组合的方法使我们能够在更高级别的抽象层面构建产品。 而不是解决整体问题，而是将其分解为组成部分，解决问题。 我们为每个问题的范围构建最佳解决方案，然后组合这些块以形成一个完整的解决方案。

## 沟通顺序过程

鉴于我们对简单，模块化，可组合的软件的信念，我们有几个原则将这些部分组合成一个连接的系统。 通信顺序过程（CSP）是一种计算模型，其中通过网络连接的自主进程能够进行通信。 我们认为，CSP方法对于管理复杂性和在面向服务架构中构建可扩展的可扩展系统是必要的。 每个服务都应被视为一个单独的进程，然后通过API与其他服务进行通信。

CSP代表着最好的方式来编写软件并组织在一起形成系统的服务。 自然本身提供了最好的例子; 甚至人体是一个相互联系的系统 - 呼吸，心血管，神经，免疫等。

## 不变性

不变性是无法改变的。 这是一个可以在很多层面上应用的概念。 不可变的最熟悉的实现是版本控制系统; 一旦代码被提交，那个提交是永远固定的。 版本控制系统，如git，广泛使用，因为它们提供了巨大的收益。 代码可以版本化，允许回滚和向前滚动。 您可以以原子方式检查和编写代码。 使用版本可以进行审核，并创建一个关于当前状态如何达到的清晰历史。 当某些事情发生时，可以使用版本历史来确定错误的起源。

不变性的概念可以扩展到基础设施的许多方面 - 应用程序源，应用程序版本和服务器状态。 我们相信使用不可变的基础架构可以提供更强大的操作，调试，版本和可视化体系。

## 通过编纂版本化

编码是相信所有进程应该被编写为代码，存储和版本化。 运营团队历来依靠口头传统，了解如何构建，升级和分类基础设施。 但信息很容易被丢失或隐藏在最需要的人身上。 关键知识的编纂促进信息共享，并防止数据丢失，因为任何进程更改都会自动存储和版本化。

HashiCorp产品都是为了遵循知识范式的编纂而设计的。 用户所做的任何更改都将进行版本化和存储，以保持过程的清晰历史。

## 通过编纂自动化

系统管理通常需要操作员手动对基础设施进行更改，使职位的职责难以扩展。 管理基础设施规模不断扩大，手工系统管理技术难以适应新规模。 自动化以更少的开销来管理更多的系统是唯一的选择。

虽然有许多自动化方法，我们促进编纂。 编码允许知识由机器执行，但操作员仍然可读。 自动化工具允许操作员提高生产率，更快地移动，并减少人为错误。 机器可以自动检测，分类和解决问题。

## 弹性系统

弹性系统可以承受意外的输入和输出。 为了实现这一点，系统必须具有期望的状态，收集关于当前状态的信息的方法以及自动调整当前状态以返回到期望状态的机制。

我们认为，将这种严格的系统应用于基础设施对于实现最高可靠性至关重要。 HashiCorp产品将通过编纂知识始终认识到所期望的状态。 他们将通过功能独立的组件收集实时信息。 他们将提供自我修复和自动恢复的工具。

## 实用主义

在遇到任何问题时，我们坚信实用主义的价值。 我们相信，不可改变，编纂，自动化和CSP的许多原则都是我们渴望的理想，而不是我们所强加的要求。 在许多情况下，实际解决方案需要重新评估我们的理想。

实用主义使我们能够评估新的想法，方法和技术，以及如何采用它们来改善HashiCorp的最佳实践。 认为这是对第一原则的妥协是一个错误，而是开放态度和谦卑，接受我们可能是错误的。 适应能力对于创新至关重要，我们为此感到自豪。
