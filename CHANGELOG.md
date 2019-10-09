v1.1.1 [2019-10-09]
-------------------

 * Add variables for extra flags to Makefile.

v1.1 [2019-07-21]
-----------------

 * Ensure that begin and end of context do not match the same substring.

v1.1rc1 [2019-07-21]
--------------------

 * Support for contexts to conditionally enable / disable deletes and replaces.
 * Option to check how often a delete / replace is applied.
 * Add script for bash completion.

v1.0.2 [2018-08-06]
-------------------

 * Do not append newline if 'append' already has one.

v1.0.1 [2018-07-01]
-------------------

 * Make sure delete doesn't remove empty lines after the pattern.

v1.0 [2018-06-24]
-----------------

Initial version that is able to replace my scripts using `sed`.

 * Deleting lines based on regular expressions.
 * Search and replace using regular expressions.
 * Appending lines at the end of files.
 * Atomic writing of configuration files, honoring permissions and ownership.
