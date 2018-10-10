<?php
require("./vendor/autoload.php");

$key = "";
$secret = "";
//$url = "https://test-idcheck.nle-tech.com/";
$url = "http://192.168.0.131:1323";

$id = "012345678912345678";
$expire = time() + 5;

//$orig = $secret . "expire=".$expire."&key=".$key."&name=".urlencode($name)."&number=".$number;
$orig = $secret . "id=".$id;
echo $orig;
echo "\n";
$sign = md5($orig);
echo $sign;
echo "\n";

$client = new GuzzleHttp\Client();
$front = fopen("./012345678912345678_front.jpg", "r");
$back = fopen("./012345678912345678_back.jpg", "r");
$res = $client->request('POST', $url . '/upload', 
	[ 
		"query" => [
			"id" => $id,
			"expire" => $expire,
			"key" => $key,
			"sign" => $sign,
		],
		"front" => $front,
		"back" => $back,
	]
);

echo $res->getBody();
echo "\n";
