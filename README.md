# Compliance Test Suite

**This project is a work in progress.**

This suite can be used to test software for compliance against specifications and to produce
reports similar to [caniuse.com](https://caniuse.com) for them.

A concrete example would be testing various database engines for compliance with SQL99.

## Goal

Good test suites allow communities and diverse software to thrive, but writing good test suites is hard.
`spec-tester` tries to automate the boring parts as much as possible so writing tests can be 
a collaboration between multiple people and machines.

## Model

Test suites are made up of three major types of components:

<dl>
<dt>Implementation</dt>
<dd>
    A single implementatation of a program, may include multiple variants e.g. versions or operating modes.
    Implementations target one or more Specifications.
</dd>

<dt>Specification</dt>
<dd>
    A single specification which is broken into multiple optional and non-optional
    compliance sections. These are usually mapped to sections of the specification.
    Sections are made up of one or more tests.
</dd>
<dt>Test</dt>
<dd>
A test for a single unit of functionality. The same test can be shared by multiple specifications
and sections. This allows specifications to be evolved over time without significant duplication.
</dd>
</dl>

