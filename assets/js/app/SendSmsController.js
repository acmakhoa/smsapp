'use strict';


smsApp.controller('SendSmsController', function($scope, $log, $http) {
  
      $scope.master = {};
      $scope.flash = {};
      $scope.send = function(sms) {
        console.log(sms);
        $scope.master = angular.copy(sms);

        var promise = $http({method: 'POST', data: sms, url: '/api/sms/create'}).
            success(function(data, status, headers, config) {
                $scope.flash = {error:0,message:"Save message successfully!"};

            }).
            error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

      };

      $scope.reset = function() {
        $scope.sms = angular.copy($scope.master);
      };

      $scope.reset();
});
