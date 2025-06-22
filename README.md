# manifold

> [man-uh-fohld] *Phonetic (Standard)*
> 
> ## adjective
> 
> 1. of many kinds; numerous and varied.
> 
> *manifold duties.*
> 
> > **Synonyms:** multitudinous, various
> > 
> > **Antonyms:** single, simple
>
> 2. having numerous different parts, elements, features, forms, etc..
> 
> > a manifold program for social reform.
> >
> > **Synonyms:** multifarious, divers, varied
> 
> 3. using, functioning with, or operating several similar or identical devices at the same time.
> 
> 4. (of paper business forms) made up of a number of sheets interleaved with carbon paper.
> 
> 5. being such or so designated for many reasons.
> 
> > a manifold enemy.
> 
> ## noun
> 
> 1. something having many different parts or features.
> 
> 2. a copy or facsimile, as of something written, such as is made by manifolding.
> 
> 3. any thin, inexpensive paper for making carbon copies on a typewriter.
> 
> 4. **Machinery.** a chamber having several outlets through which a liquid or gas is distributed or gathered.
> 
> 5. **Philosophy.** (in Kantian epistemology) the totality of discrete items of experience as presented to the mind; the constituents of a sensory experience.
> 
> 6. **Mathematics.** a topological space that is connected and locally Euclidean.
> 
> ## verb (used with object)
> 
> to make copies of, as with carbon paper.
> 
> *from [dictionary.com](https://www.dictionary.com/browse/manifold)*

Manifold is an event bus type protocol that combines post-office/repository access with publish/subscribe 
distribution of messages.

The name relates to the way that it is designed to be so simple that it can perform many different kinds of messaging 
and data distribution, replication and storage for any purpose. specifically it's goal is to function in place of 
standard web services, as well as enabling the necessary real-time push of new data across the network, while allowing 
both active and passive methods of synchronisation between nodes in the network.

It is designed to be as simple as possible, for example, using a simple sentinel based encoding for events and filters 
that only require escaping newlines and backslashes, instead of a complex and varied escape scheme like found with JSON.

This is important not only for easy implementation but also because cryptographic signatures are only valid on an exact
string of bytes, and Manifold events are identified by their SHA256 hash. 

The ordering of fields is rigid and requires exact ordering, and instead of adding a "kind" to specify the types, 
the tags serve the purpose of marking the content encoding, purpose, application type, as well as any other use such 
as event and user pubkey references for threaded discussions.