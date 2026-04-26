package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var approvalCmd = &cobra.Command{
	Use:   "approval",
	Short: "Manage change-approval requests for secret paths",
}

var requestApprovalCmd = &cobra.Command{
	Use:   "request <path> <version> <requestor>",
	Short: "Create a pending approval request for a secret version",
	Args:  cobra.ExactArgs(3),
	RunE:  runRequestApproval,
}

var reviewApprovalCmd = &cobra.Command{
	Use:   "review <path> <reviewer> <approve|reject>",
	Short: "Approve or reject a pending approval request",
	Args:  cobra.ExactArgs(3),
	RunE:  runReviewApproval,
}

var getApprovalCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Show the current approval record for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetApproval,
}

var approvalReason string

func init() {
	requestApprovalCmd.Flags().StringVar(&approvalReason, "reason", "", "Reason for the change request")
	approvalCmd.AddCommand(requestApprovalCmd, reviewApprovalCmd, getApprovalCmd)
	rootCmd.AddCommand(approvalCmd)
}

func runRequestApproval(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}
	requestor := args[2]

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	req, err := vault.RequestApproval(client, path, version, requestor, approvalReason)
	if err != nil {
		return err
	}

	fmt.Printf("Approval request created\n")
	fmt.Printf("  Path:      %s\n", req.Path)
	fmt.Printf("  Version:   %d\n", req.Version)
	fmt.Printf("  Requestor: %s\n", req.Requestor)
	fmt.Printf("  Status:    %s\n", req.Status)
	if req.Reason != "" {
		fmt.Printf("  Reason:    %s\n", req.Reason)
	}
	return nil
}

func runReviewApproval(cmd *cobra.Command, args []string) error {
	path := args[0]
	reviewer := args[1]
	decision := args[2]

	var approved bool
	switch decision {
	case "approve":
		approved = true
	case "reject":
		approved = false
	default:
		return fmt.Errorf("decision must be 'approve' or 'reject', got %q", decision)
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	req, err := vault.ReviewApproval(client, path, reviewer, approved)
	if err != nil {
		return err
	}

	fmt.Printf("Approval updated: %s (reviewed by %s)\n", req.Status, req.Reviewer)
	return nil
}

func runGetApproval(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	req, err := vault.GetApproval(client, args[0])
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Path:\t%s\n", req.Path)
	fmt.Fprintf(w, "Version:\t%d\n", req.Version)
	fmt.Fprintf(w, "Requestor:\t%s\n", req.Requestor)
	fmt.Fprintf(w, "Reason:\t%s\n", req.Reason)
	fmt.Fprintf(w, "Status:\t%s\n", req.Status)
	if req.Reviewer != "" {
		fmt.Fprintf(w, "Reviewer:\t%s\n", req.Reviewer)
	}
	fmt.Fprintf(w, "Created:\t%s\n", req.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Updated:\t%s\n", req.UpdatedAt.Format("2006-01-02 15:04:05 UTC"))
	return w.Flush()
}
