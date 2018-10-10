<?php
require("./vendor/autoload.php");
$key = "";
$secret = "";
$url = "http://192.168.0.131:1323";
$id = "012345678912345678";
$expire = time() + 5;
$orig = $secret . "expire=".$expire ."&id=".$id . "&key=".$key;
echo $orig;
echo "\n";
$sign = md5($orig);
echo $sign;
echo "\n";
$client = new GuzzleHttp\Client();
$front = fopen("/root/idcheck/012345678912345678_front.jpg", "r");
$back = fopen("/root/idcheck/012345678912345678_back.jpg", "r");
$res = $client->request('POST', $url . '/upload',
  [
    "query" => [
      "id" => $id,
      "expire" => $expire,
      "key" => $key,
      "sign" => $sign,
    ],
    'multipart' => [
      [
        'name' => 'front',
        'contents' => $front,
        'filename' => '/root/idcheck/012345678912345678_front.jpg',
      ],
      [
        'name' => 'back',
        'contents' => $back,
        'filename' => '/root/idcheck/012345678912345678_back.jpg',
      ],
    ]
  ]
);
echo $res->getBody();
echo "\n"
