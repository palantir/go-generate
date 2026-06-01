package fail

fail

/*
This is a non-compiling file that has been added to explicitly ensure that CI fails.
It also contains the command that caused the failure and its output.
Remove this file if debugging locally.

./godelw verify failed after updating godel plugins and assets

Command that caused error:
./godelw verify

Output:
Error: failed to resolve 1 plugin(s):
    failed to resolve artifact com.palantir.godel-mod-plugin:mod-plugin:1.52.0 using resolvers:
    failed to resolve artifact at https://github.com/palantir/godel-mod-plugin/releases/download/v1.52.0/mod-plugin-1.52.0-linux-amd64.tgz: request for URL https://github.com/palantir/godel-mod-plugin/releases/download/v1.52.0/mod-plugin-1.52.0-linux-amd64.tgz returned status code 404

*/
