At last we should support the following features
    - Enum enumerations (we can craft helper methods)
    - Tuples (we can craft helper methods only and i think that's enough)
    - Multiline hash tag support: useful for implementing enumerations or providing multiple options to macro.
    - BTW macros should accept input: to implement former feature

So at the time we need to pass `*ast.GenDecl` to macro definition.

*** Also there should be an option to write the output code in a custom path!

----------------------------------------------------------------------------------------------

Applications:
 - Add `Clone` derive macro to enable deep clone for every types that need to deep clone feature
 	every fields of a struct marked as `Clone` must be `Clone` too. underlying type of a new type also need to be `Clone`
 - JSON efficient serialization/deserialization.
 - Implement efficient and concise SQL model bindings
