'use strict';


smsApp.controller('ListDetailController', function($scope, $log, $http) {
  
    $scope.sms = [];
    $scope.list = [];
    $scope.loading = true;
    $scope.listId = listId;
    var filter = $scope.filter = {
       
    };
    var alllistSms = [];
    var listItem;
    $scope.delete = function(sms,idx) {            
        var promise = $http({method: 'POST', data: null, url: '/api/list/'+$scope.listId+'/delete/'+sms.Id}).
            success(function(data, status, headers, config) {
                alllistSms.splice(idx, 1);

            }).
            error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

      };

    $scope.send = function(sms,idx) {            
        var promise = $http({method: 'POST', data: null, url: '/api/list/'+$scope.listId+'/send/'+sms.Id}).
            success(function(data, status, headers, config) {
               alert("Send success");
            }).
            error(function(data, status, headers, config) {
                alert("Send error");
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

      };


    var promise = $http({method: 'GET', data: null, url: '/api/list/'+$scope.listId}).
            success(function(res, status, headers, config) {
                listItem=res.data.list;
                for (var i = 0; i < res.data.sms.length; i++) {
                    alllistSms.push(res.data.sms[i]);
                }

            }).
            error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        promise.then(function() {
            $scope.loading = false;
            $scope.sms = alllistSms;
            $scope.list = listItem;

            $scope.$watch('filter', filterAndSortListSms, true);
        });

        $scope.$watch('filter', filterAndSortListSms, true);
        function filterAndSortListSms() {
        
             //sort
            $scope.sms.sort(function(a, b) {

                if (a.Created > b.Created) {
                    return filter.sortAsc ? 1 : -1;
                }

                if (a.Created < b.Created) {
                    return filter.sortAsc ? -1 : 1;
                }

                return 0;
            });
        }
    }


);