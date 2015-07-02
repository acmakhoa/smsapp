'use strict';


smsApp.controller('SettingController', function($scope, $log, $http) {
  
      $scope.modems = [];
      $scope.flash = {};
      $scope.rescan = function() {                
        var promise = $http({method: 'GET', data: {}, url: '/api/device/rescan'}).
            success(function(data, status, headers, config) {
              if(data.status==200){
                $scope.modems=data.data.modems;
                console.log(data.data);
                $scope.flash = {error:0,message:"Save message successfully!"};  
              }
              
            }).
            error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

      };

      $scope.select = function() {                
        var promise = $http({method: 'POST', data: {name:$scope.name}, url: '/api/device/select'}).
            success(function(data, status, headers, config) {
              if(data.status==200){                
                $scope.flash = {error:0,message:"Save message successfully!"};  
              }
              
            }).
            error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

      };

      



      var promise = $http({method: 'GET', data: {}, url: '/api/device/rescan'}).
      success(function(data, status, headers, config) {
        if(data.status==200){
          $scope.modems=data.data.modems;
          console.log(data.data);
          $scope.flash = {error:0,message:"Save message successfully!"};  
        }

      }).
      error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
              });

      

});
