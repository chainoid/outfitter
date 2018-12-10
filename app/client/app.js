// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function ($scope, appFactory) {

	// Intendant page
	$("#all_groups").hide();
	$("#all_users").hide();

	$("#error_add_group").hide();
	$("#success_add_group").hide();

	$("#error_add_user").hide();
	$("#success_add_user").hide();

	$("#success_generated").hide();
	$("#error_generated").hide();

	//

	$("#error_query").hide();
	$("#error_sender").hide();
	$("#error_query_id").hide();
	$("#error_query_student").hide();
	$("#error_prepare_delivery").hide();
	$("#error_pass_exam").hide();
	$("#error_student_record").hide();
	$("#item_list").hide();
	
	$("#error_exam_source").hide();
	$("#error_old_exam").hide();
	$("#success_exam").hide();

	// Fighter page	
	$("#error_user_record").hide();
	$("#user_record").hide();
	$("#user_record2").hide();

	

	$("#take_form").hide();

	$scope.queryAllGroups = function () {

		appFactory.queryAllGroups(function (data) {
			var array = [];
			for (var i = 0; i < data.length; i++) {
				data[i].Record.Key = data[i].Key;
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return a.groupName.localeCompare(b.groupName);
			});
			$scope.all_groups = array;
			$("#all_groups").show();
		});
	}

	$scope.addGroup = function () {

		appFactory.addGroup($scope.newGroup, function (data) {

			if (data == "Could not locate unpassed test") {
				$("#error_add_group").show();
				$("#success_add_group").hide();
			} else {
				$("#error_add_group").hide();
				$("#success_add_group").show();
			}

			$scope.exam_result = data;
		});
	}

    $scope.addUser = function () {

		appFactory.addUser($scope.user, function (data) {

			if (data == "Could not locate unpassed test") {
				$("#error_add_user").show();
				$("#success_add_user").hide();
			} else {
				$("#error_add_user").hide();
				$("#success_add_user").show();
			}

			$scope.exam_result = data;
		});
	}

	$scope.queryAllUsers = function () {

		appFactory.queryAllUsers(function (data) {
			var array = [];
			for (var i = 0; i < data.length; i++) {
				data[i].Record.Key = data[i].Key;
				array.push(data[i].Record);
			}
			array.sort(function (a, b) {
				return a.groupName.localeCompare(b.groupName);
			});
			$scope.all_users = array;
			$("#all_users").show();
		});
	}

	$scope.generateSetForGroup = function () {

		appFactory.generateSetForGroup($scope.generator, function (data) {
			$scope.generated_set_for_group = data;

			if ($scope.generated_set_for_group == "error_generated") {
				console.log()
				$("#error_generated").show();
			} else {
				$("#error_generated").hide();
				$("#success_generated").show();
			}

		});
	}

	$scope.getUserRecord = function () {
		
		var id = $scope.id;

		appFactory.getUserRecord(id, function(data){

			$scope.user_record = data;

			if ($scope.user_record == "User record not found"){
				console.log()
				$("#error_user_record").show();
				$("#user_record").hide();
				$("#user_record2").hide();
				
			} else{
				$("#error_user_record").hide();
				$("#user_record").show();
				$("#user_record2").show();
			}
		});
	}

	$scope.prepareForDelivery = function () {

		var order = $scope.order;

		appFactory.prepareForDelivery(order, function (data) {

			if (data == "No group/item found") {
				console.log("No group/item found");
				$("#error_prepare_delivery").show();
				$("#item_list").hide();
			
			} else {
				$("#error_prepare_delivery").hide();
				$("#item_list").show();
				$("#take_form").hide(); 
			}

			var array = [];
			for (var i = 0; i < data.length; i++) {
				data[i].Record.Key = data[i].Key;
				array.push(data[i].Record);
			}
			array.sort(function (a, b) {
				return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.item_list = array;
		});
	}


	$scope.beforeDeliveryItem = function (item) {
		        
          if (item.rate != "") {
			$("#takeTheTestId").hide();	 
			$("#take_form").hide(); 
		  } else {
			$("#takeTheTestId").show();	
			$("#take_form").show();
			$("#success_exam").hide();
		  }
		  $scope.examcase = item;
	}

	$scope.deliveryItem = function () {

		var examcase = $scope.examcase;

		appFactory.deliveryItem(examcase, function (data) {

			if (data == "Could not locate unpassed test") {
				$("#error_exam_source").show();
				$("#success_exam").hide();
			} else {	
				$("#error_exam_source").hide();
				$("#success_exam").show();
			}
			
			$scope.exam_result = data;
		});
	}

});


// Angular Factory
app.factory('appFactory', function ($http) {

	var factory = {};

	factory.queryAllGroups = function (callback) {

		$http.get('/get_all_groups/').success(function (output) {
			callback(output)
		});
	}


	factory.addGroup = function (data, callback) {

		var newGroup =  data.groupName + "-" + data.description;

		$http.get('/add_group/' + newGroup).success(function (output) {
			callback(output)
		});
	}


	factory.addUser = function (data, callback) {

		var student = data.studentId + "-" + data.studentName + "-" + data.groupName + "-" + data.description;

		$http.get('/add_user/' + student).success(function (output) {
			callback(output)
		});
	}

	factory.queryAllUsers = function (callback) {

		$http.get('/query_all_users/').success(function (output) {
			callback(output)
		});
	}
	
	factory.generateSetForGroup = function (generator, callback) {

		var generator = generator.groupName + "-" + generator.courseName + "-" + generator.teacherName;

		$http.get('/generate_set_for_group/' + generator).success(function (output) {
			callback(output)
		});
	}

	factory.getUserRecord = function (id, callback) {
		$http.get('/get_user_record/' + id).success(function (output) {
			callback(output)
		});
	}

	factory.prepareForDelivery = function (exam, callback) {

		var params = exam.group + "-" + exam.course;

		$http.get('/prepare_for_delivery/' + params).success(function (output) {
			callback(output)
		});
	}

	factory.deliveryItem = function (input, callback) {

		var params = input.studentId + "-" + input.course + "-" + input.rate;

		$http.get('/delivery_item/' + params).success(function (output) {
			callback(output)
		});
	}

	return factory;
});
