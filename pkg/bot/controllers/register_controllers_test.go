package controllers

import "testing"

func TestFacultyController_Verify(t *testing.T) {
	fc := &FacultyController{}
	fc.Verify(nil)
}
