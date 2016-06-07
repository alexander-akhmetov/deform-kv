deformkv
==============

Simple key-value wrapper for deform.io

deform.io [quickstart](http://deformio.github.io/docs/quickstart/).

You have to create project:

    deform project create -d '{"_id": "project_name", "name": "My key-value storage"}'

And [token](http://deformio.github.io/docs/quickstart/#creating-a-token)

How to use:

    deform := deformkv.NewClient(project, collection, token)
    deform.Set("my-key", "some-value")
    value, err := deform.Get("my-key")  // returns "some-value"
