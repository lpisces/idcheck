<?php
require("./vendor/autoload.php");

$key = "testkey";
$secret = "testsecret";
//$url = "https://idcheck.nle-tech.com";
$url = "http://192.168.0.131:1323";

//$sms_content = urlencode("测试");
$sms_content ="测试";
$sms_sign = "测试";
$sms_mobile = "18674469015";
$expire = time() + 5;

$orig = sprintf("%sexpire=%s&key=%s&sms_content=%s&sms_mobile=%s&sms_sign=%s", $secret, $expire, $key, urlencode($sms_content), $sms_mobile, urlencode($sms_sign));
echo $orig;
echo "\n";
$sign = md5($orig);
echo $sign;
echo "\n";

$client = new GuzzleHttp\Client();
$res = $client->request('POST', $url . '/sms',
  [
    "query" => [
      "sms_content" => $sms_content,
      "sms_mobile" => $sms_mobile,
      "sms_sign" => $sms_sign,
      "key" => $key,
      "sign" => $sign,
      "expire" => $expire,
    ],
  ]
);

echo $res->getBody();
echo "\n";
