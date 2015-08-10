curl 'localhost:8011/case?Status="tested"&Page=0'

curl localhost:8011/2b0bc394cdc13a748842ff86e058a761 > a.tar.gz
tar xzvf a.tar.gz

curl localhost:8011/2b0bc394cdc13a748842ff86e058a761 > a.report
cat a.report

curl -d {} localhost:8011/case
