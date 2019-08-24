- First commit
  Worked out the bases for the handling of the keys files.
  The keys file will have permission `0400`
  This means only the user who run the golang app will have access to it...and root of course.
  Tests cover 100% of the functions, although they are all grouped in a single method.