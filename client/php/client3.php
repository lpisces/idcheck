<?php
require("./vendor/autoload.php");
$key = "";
$secret = "";
$url = "http://192.168.0.131:1323";
$ids = "012345678912345678";
$expire = time() + 5;
$orig = $secret . "expire=".$expire . "&key=".$key;
echo $orig;
echo "\n";
$sign = md5($orig);
echo $sign;
echo "\n";
$client = new GuzzleHttp\Client();
$res = $client->request('POST', $url . '/download',
  [
    "query" => [
      "expire" => $expire,
      "key" => $key,
      "sign" => $sign,
    ],
    "form_params" => [
      "ids" => $ids
    ],
  ]
);
echo $res->getBody();
echo "\n";
