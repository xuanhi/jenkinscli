package jenkins

import "fmt"

const a string = `<!DOCTYPE html>
<html>
	<style>
	html,
	body,
	.box .content {
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: pink;
	}
	
	.box {
		width: 300px;
		height: 300px;
		box-sizing: border-box;
		padding: 15px;
		position: relative;
		overflow: hidden;
	}
	
	.box::before {
		content: '';
		position: absolute;
		width: 150%;
		height: 150%;
		background: repeating-linear-gradient(
				white 0%,
				white 7.5px,
				hotpink 7.5px,
				hotpink 15px,
				white 15px,
				white 22.5px,
				hotpink 22.5px,
				hotpink 30px);
		transform: translateX(-20%) translateY(-20%) rotate(-45deg);
		animation: animate 20s linear infinite;
	}
	
	.box .content {
		position: relative;
		background-color: white;
		flex-direction: column;
		box-sizing: border-box;
		padding: 30px;
		text-align: center;
		font-family: sans-serif;
		z-index: 2;
	}
	
	.box,
	.box .content {
		box-shadow: 0 0 2px deeppink,
					0 0 5px rgba(0, 0, 0, 1),
					inset 0 0 5px rgba(0, 0, 0, 1);
		border-radius: 10px;
	}
	
	.box .content h2 {
		color: deeppink;
	}
	
	.box .content p {
		color: dimgray;
	}
	
	@keyframes animate {
		from {
			background-position: 0;
		}
	
		to {
			background-position: 0 450px;
		}
	}
	</style>
	<head>
		<meta charset="utf-8" />
		<title></title>
	</head>
	<body>
		<div class="box">
		  <div class="content">
		    <h2>Jenkins有新构建啦</h2>
		    <p>流水线：`
const b string = `</p>
<p>构建号:`
const c string = `</p>
<p>构建结果:`
const d string = `</p>
</div>
</div>
</body>
</html>`

func Emailpost(number int64, name, result string) string {
	return fmt.Sprintln(a, name, b, number, c, result, d)
}
