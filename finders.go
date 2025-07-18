package main

import (
	"context"
)

// findDirectoryAuthor finds the most recent author in the affected directories.
func (rf *ReviewerFinder) findDirectoryAuthor(ctx context.Context, pr *PullRequest, files []string) string {
	directories := rf.getDirectories(files)
	
	for _, dir := range directories {
		if author := rf.findRecentAuthorInDirectory(ctx, pr.Owner, pr.Repository, dir); author != "" {
			if author != pr.Author && rf.isValidReviewer(ctx, pr, author) {
				return author
			}
		}
	}
	
	return ""
}

// findProjectAuthor finds the most recent author in the project.
func (rf *ReviewerFinder) findProjectAuthor(ctx context.Context, pr *PullRequest) string {
	authors := rf.findRecentAuthorsInProject(ctx, pr.Owner, pr.Repository, 10)
	
	for _, author := range authors {
		if author != pr.Author && rf.isValidReviewer(ctx, pr, author) {
			return author
		}
	}
	
	return ""
}

// findCommenterReviewer finds recent PR commenters who can review.
func (rf *ReviewerFinder) findCommenterReviewer(ctx context.Context, pr *PullRequest, excludeAuthor string) string {
	excludeUsers := []string{pr.Author}
	if excludeAuthor != "" {
		excludeUsers = append(excludeUsers, excludeAuthor)
	}
	
	commenters, err := rf.getRecentPRCommenters(ctx, pr.Owner, pr.Repository, excludeUsers)
	if err != nil {
		return ""
	}
	
	for _, commenter := range commenters {
		if rf.isValidReviewer(ctx, pr, commenter) {
			return commenter
		}
	}
	
	return ""
}

// findDirectoryReviewer finds the most active reviewer in the affected directories.
func (rf *ReviewerFinder) findDirectoryReviewer(ctx context.Context, pr *PullRequest, files []string, excludeAuthor string) string {
	directories := rf.getDirectories(files)
	
	for _, dir := range directories {
		if reviewer := rf.findActiveReviewerInDirectory(ctx, pr.Owner, pr.Repository, dir); reviewer != "" {
			if reviewer != pr.Author && reviewer != excludeAuthor && rf.isValidReviewer(ctx, pr, reviewer) {
				return reviewer
			}
		}
	}
	
	return ""
}

// findProjectReviewer finds the most active reviewer in the project.
func (rf *ReviewerFinder) findProjectReviewer(ctx context.Context, pr *PullRequest, excludeAuthor string) string {
	reviewers := rf.findActiveReviewersInProject(ctx, pr.Owner, pr.Repository, 10)
	
	for _, reviewer := range reviewers {
		if reviewer != pr.Author && reviewer != excludeAuthor && rf.isValidReviewer(ctx, pr, reviewer) {
			return reviewer
		}
	}
	
	return ""
}

// findRecentAuthorInDirectory finds the most recent commit author in a directory.
func (rf *ReviewerFinder) findRecentAuthorInDirectory(ctx context.Context, owner, repo, directory string) string {
	prs, err := rf.getRecentPRsInDirectory(ctx, owner, repo, directory)
	if err != nil {
		return ""
	}
	
	fc := make(frequencyCounter)
	for _, pr := range prs {
		fc.add(pr.Author)
	}
	return fc.best()
}

// findActiveReviewerInDirectory finds the most active reviewer in a directory.
func (rf *ReviewerFinder) findActiveReviewerInDirectory(ctx context.Context, owner, repo, directory string) string {
	prs, err := rf.getRecentPRsInDirectory(ctx, owner, repo, directory)
	if err != nil {
		return ""
	}
	
	fc := make(frequencyCounter)
	for _, pr := range prs {
		fc.addAll(pr.Reviewers)
	}
	return fc.best()
}

// findRecentAuthorsInProject finds recent commit authors in the project.
func (rf *ReviewerFinder) findRecentAuthorsInProject(ctx context.Context, owner, repo string, limit int) []string {
	prs, err := rf.getRecentPRsInProject(ctx, owner, repo)
	if err != nil {
		return nil
	}
	
	fc := make(frequencyCounter)
	for _, pr := range prs {
		fc.add(pr.Author)
	}
	return fc.top(limit)
}

// findActiveReviewersInProject finds active reviewers in the project.
func (rf *ReviewerFinder) findActiveReviewersInProject(ctx context.Context, owner, repo string, limit int) []string {
	prs, err := rf.getRecentPRsInProject(ctx, owner, repo)
	if err != nil {
		return nil
	}
	
	fc := make(frequencyCounter)
	for _, pr := range prs {
		fc.addAll(pr.Reviewers)
	}
	return fc.top(limit)
}

