ALTER TABLE "contents" RENAME TO "artifacts";--> statement-breakpoint
ALTER TABLE "_labelsToContent" RENAME COLUMN "content_id" TO "artifact_id";--> statement-breakpoint
ALTER TABLE "_runToContent" RENAME COLUMN "content_id" TO "artifact_id";--> statement-breakpoint
ALTER TABLE "_labelsToContent" DROP CONSTRAINT "_labelsToContent_content_id_contents_id_fk";
--> statement-breakpoint
ALTER TABLE "_parentToChild" DROP CONSTRAINT "_parentToChild_child_id_contents_id_fk";
--> statement-breakpoint
ALTER TABLE "_parentToChild" DROP CONSTRAINT "_parentToChild_parent_id_contents_id_fk";
--> statement-breakpoint
ALTER TABLE "_runToContent" DROP CONSTRAINT "_runToContent_content_id_contents_id_fk";
--> statement-breakpoint
ALTER TABLE "artifacts" DROP CONSTRAINT "contents_orgname_organizations_id_fk";
--> statement-breakpoint
ALTER TABLE "artifacts" DROP CONSTRAINT "contents_parent_id_contents_id_fk";
--> statement-breakpoint
ALTER TABLE "artifacts" DROP CONSTRAINT "contents_producer_id_runs_id_fk";
--> statement-breakpoint
ALTER TABLE "_labelsToContent" DROP CONSTRAINT "_labelsToContent_label_id_content_id_pk";--> statement-breakpoint
ALTER TABLE "_runToContent" DROP CONSTRAINT "_runToContent_run_id_content_id_pk";--> statement-breakpoint
ALTER TABLE "_labelsToContent" ADD CONSTRAINT "_labelsToContent_label_id_artifact_id_pk" PRIMARY KEY("label_id","artifact_id");--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_run_id_artifact_id_pk" PRIMARY KEY("run_id","artifact_id");--> statement-breakpoint
ALTER TABLE "_labelsToContent" ADD CONSTRAINT "_labelsToContent_artifact_id_artifacts_id_fk" FOREIGN KEY ("artifact_id") REFERENCES "public"."artifacts"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_parentToChild" ADD CONSTRAINT "_parentToChild_child_id_artifacts_id_fk" FOREIGN KEY ("child_id") REFERENCES "public"."artifacts"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_parentToChild" ADD CONSTRAINT "_parentToChild_parent_id_artifacts_id_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."artifacts"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_artifact_id_artifacts_id_fk" FOREIGN KEY ("artifact_id") REFERENCES "public"."artifacts"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_parent_id_artifacts_id_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."artifacts"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_producer_id_runs_id_fk" FOREIGN KEY ("producer_id") REFERENCES "public"."runs"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;