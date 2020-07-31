// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"net/http"

	"github.com/google/exposure-notifications-verification-server/pkg/controller"
	"github.com/google/exposure-notifications-verification-server/pkg/controller/flash"
	"github.com/gorilla/mux"
)

func (c *Controller) HandleDelete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		flash := flash.FromContext(w, r)

		user := controller.UserFromContext(ctx)
		if user == nil {
			flash.Error("Unauthorized.")
			http.Redirect(w, r, "/signout", http.StatusSeeOther)
			return
		}

		realm := controller.RealmFromContext(ctx)
		if realm == nil {
			flash.Error("Select a realm to continue.")
			http.Redirect(w, r, "/realm", http.StatusSeeOther)
			return
		}

		// TODO(sethvargo): switch to a form post parameter instead - this leaks
		// emails in request logs.
		vars := mux.Vars(r)
		email := vars["email"]

		user, err := c.db.FindUser(email)
		if err != nil {
			flash.Error("Failed to find user: %v", err)
			http.Redirect(w, r, "/users", http.StatusSeeOther)
			return
		}

		if err := realm.DeleteUserFromRealm(c.db, user); err != nil {
			flash.Error("Failed to delete user: %v", err)
			http.Redirect(w, r, "/users", http.StatusSeeOther)
			return
		}

		flash.Alert("Deleted User %v", user.Email)
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	})
}